package main

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/bmatcuk/doublestar/v2"
	"github.com/cheggaaa/pb/v3"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zip"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var VERBOSE bool
var DEBUG bool

type Unit struct {
	name             string
	sources          []string
	destination      string
	excludes         []string
	archiveType      string
	addSubfolder     bool
	enabled          bool
	useAbsolutePaths bool
}

type Config struct {
	units []Unit
}

type BackupFileMetadata struct {
	path           string
	backupBasePath string
}

func (config Config) FromYaml(yamlData []byte) (Config, error) {
	// Create a config object from yaml byte array
	log.Println("Parsing config yaml")

	// Helper struct for parsing the yaml
	type YamlUnit struct {
		Sources          *[]string `yaml:"sources"`
		Destination      *string   `yaml:"destination"`
		Excludes         *[]string `yaml:"excludes"`
		ArchiveType      *string   `yaml:"archive_type"`
		AddSubfolder     *bool     `yaml:"add_subfolder"`
		Enabled          *bool     `yaml:"enabled"`
		UseAbsolutePaths *bool     `yaml:"use_absolute_paths"`
	}

	// Sample yaml config:
	/*
		backup_unit_name:
		  sources:
		    - C:\Users\admin\Documents\
		    - 'C:\Users\admin\Dropbox'
		    - "C:\\Users\\admin\\AppData\\Roaming\\.minecraft"
		  destination: 'C:\backups'
		  excludes:
		    - "*.zip"
		    - "*.rar"
		  archive_type: "tar.gz"
		  add_subfolder: false
	*/

	unitMap := make(map[string]YamlUnit)

	unmarshalErr := yaml.Unmarshal(yamlData, &unitMap)
	if unmarshalErr != nil {
		log.Fatalf("Unmarshal error: %v", unmarshalErr)
	}

	// After parsing the yaml into unitMap, we iterate over all available units
	for unitName, yamlUnit := range unitMap {
		unit := Unit{}

		// Set defaults
		unit.enabled = true
		if yamlUnit.Enabled != nil {
			unit.enabled = *yamlUnit.Enabled
		}

		unit.addSubfolder = false
		if yamlUnit.AddSubfolder != nil {
			unit.addSubfolder = *yamlUnit.AddSubfolder
		}

		unit.archiveType = "tar.gz"
		if yamlUnit.ArchiveType != nil {
			unit.archiveType = *yamlUnit.ArchiveType
		}

		unit.excludes = []string{}
		if yamlUnit.Excludes != nil {
			unit.excludes = *yamlUnit.Excludes
		}

		unit.useAbsolutePaths = true
		if yamlUnit.UseAbsolutePaths != nil {
			unit.useAbsolutePaths = *yamlUnit.UseAbsolutePaths
		}

		if yamlUnit.Sources == nil || yamlUnit.Destination == nil {
			log.Fatalf("Sources or destination can't be parsed for unit '%s'", unitName)
		} else {
			unit.sources = *yamlUnit.Sources
			unit.destination = *yamlUnit.Destination
		}
		unit.name = unitName

		config.units = append(config.units, unit)
	}

	return config, nil
}

func handleExclude(filePath string, excludePattern string) bool {
	if excludePattern == "" {
		return false
	}

	filePath = filepath.ToSlash(filePath)
	if !strings.Contains(excludePattern, "/") {
		// When there is no forward slash in the pattern, we want to match a file
		lastIndex := strings.LastIndex(filePath, "/")
		filePath = filePath[lastIndex+1:]
	}

	matched, matchErr := doublestar.Match(excludePattern, filePath)

	if matched && DEBUG {
		log.Printf("Excluding path '%s' because pattern '%s' matched", filePath, excludePattern)
	}

	if matchErr != nil {
		log.Println(matchErr)
	}

	return matched
}

func handleExcludes(path string, excludes []string) bool {
	// Checks if the path is excluded by any of the given exclude patterns
	for _, excludePattern := range excludes {
		matched := handleExclude(path, excludePattern)

		if matched {
			return true
		}
	}
	return false
}

func getFiles(sourcePath string, excludes []string) ([]BackupFileMetadata, error) {
	// Returns all file paths recursively within a certain source directory
	var pathsToBackup []BackupFileMetadata

	_, statErr := os.Stat(sourcePath)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			log.Printf("Source directory '%s' does not exist.\n", sourcePath)
		}
		return nil, statErr
	}

	// Recursively check directories for files. Add all that do not match the exclusion filters
	err := filepath.Walk(sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			isExcluded := handleExcludes(path, excludes)
			if isExcluded {
				return nil
			}

			fileMetadata := BackupFileMetadata{
				path:           path,
				backupBasePath: sourcePath,
			}
			pathsToBackup = append(pathsToBackup, fileMetadata)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return pathsToBackup, err
}

func validatePath(path string, mustBeDir bool) bool {
	// Checks if a file/directory exists
	file, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("File '%s' does not exist.", path)
		}
		return false
	}
	if mustBeDir {
		return file.IsDir()
	}
	return true
}

func getPathInArchive(filePath string, backupBasePath string, unit Unit) string {
	// Remove the base path from the file path within the archive, if option is set
	pathInArchive := filePath
	if !unit.useAbsolutePaths {
		parentBasePath := filepath.Dir(backupBasePath)
		pathInArchive = strings.ReplaceAll(filePath, parentBasePath, "")

		if strings.HasPrefix(pathInArchive, "\\") {
			pathInArchive = pathInArchive[1:]
		}
	}
	return pathInArchive
}

func writeArchive(backupArchivePath string, filesToBackup []BackupFileMetadata, unit Unit) {
	archiveFile, err := os.Create(backupArchivePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer archiveFile.Close()

	switch unit.archiveType {
	case "tar.gz":
		writeTar(archiveFile, filesToBackup, unit)
	case "zip":
		writeZip(archiveFile, filesToBackup, unit)
	default:
		log.Fatalf("Can't handle archive type '%s'", unit.archiveType)
	}
}

func writeTar(archiveFile *os.File, filesToBackup []BackupFileMetadata, unit Unit) {
	// set up the gzip and tar writer
	gw := gzip.NewWriter(archiveFile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Init progress bar
	bar := pb.StartNew(len(filesToBackup))
	for i := range filesToBackup {
		fileMetadata := filesToBackup[i]
		filePath := fileMetadata.path

		pathInArchive := getPathInArchive(filePath, fileMetadata.backupBasePath, unit)

		if err := addFileToTar(tw, filePath, pathInArchive); err != nil {
			log.Println(err)
		}
		bar.Increment()
	}
	bar.Finish()
}

func writeZip(archiveFile *os.File, filesToBackup []BackupFileMetadata, unit Unit) {
	zw := zip.NewWriter(archiveFile)
	defer zw.Close()

	bar := pb.StartNew(len(filesToBackup))
	for i := range filesToBackup {
		fileMetadata := filesToBackup[i]
		filePath := fileMetadata.path

		pathInArchive := getPathInArchive(filePath, fileMetadata.backupBasePath, unit)

		if err := addFileToZip(zw, filePath, pathInArchive); err != nil {
			log.Println(err)
		}
		bar.Increment()
	}
	bar.Finish()
}

func writeBackup(filesToBackup []BackupFileMetadata, unit Unit) {
	now := time.Now()
	timeStamp := now.Format("2006-01-02_15-04")
	backupBasePath := unit.destination

	if unit.addSubfolder {
		newBackupBasePath := filepath.Join(unit.destination, unit.name)
		pathExists := validatePath(newBackupBasePath, true)

		if !pathExists {
			log.Printf("Backup path '%s' does not exist.\n", newBackupBasePath)
			mkdirErr := os.Mkdir(newBackupBasePath, 0777)
			if mkdirErr != nil {
				log.Fatalf("Can't create backup directory '%s'", newBackupBasePath)
			}
		}
		backupBasePath = newBackupBasePath
	}

	backupArchiveName := unit.name + "-" + timeStamp + "." + unit.archiveType
	backupArchivePath := filepath.Join(backupBasePath, backupArchiveName)

	// TODO check if archive already exists. If yes, append -1 to it and try again

	writeArchive(backupArchivePath, filesToBackup, unit)
	log.Printf("Archive created successfully at '%s'", backupArchivePath)
}

func backupUnit(unit Unit) {
	// Start backup for a single unit. Each backup creates a single archive file
	if !unit.enabled {
		log.Printf("Skipping backup for unit '%s' because it's disabled.\n", unit.name)
		return
	}

	log.Printf("Creating backup for unit '%s'\n", unit.name)
	var filesToBackup []BackupFileMetadata
	var processedSources []string

	// Check all source files from the disk in the specified source directories
	for _, sourcePath := range unit.sources {
		sourcePath = filepath.Clean(sourcePath)

		// Prevent duplicate source paths
		for _, processedPath := range processedSources {
			if sourcePath == processedPath {
				log.Printf("Found duplicate source path '%s'. Skipping!", sourcePath)
				continue
			}
		}
		processedSources = append(processedSources, sourcePath)

		files, err := getFiles(sourcePath, unit.excludes)
		if err != nil {
			log.Printf("Error for unit '%s' while reading directory '%s'! Skipping!", unit.name, sourcePath)
			continue
		}
		filesToBackup = append(filesToBackup, files...)
	}

	if len(filesToBackup) == 0 {
		log.Printf("No files found for sources in unit '%s'. Creating no backup!", unit.name)
		return
	}
	writeBackup(filesToBackup, unit)
}

func runBackup(config Config) {
	// Start the backup(s) defined in the config object
	for _, unit := range config.units {
		backupUnit(unit)
	}
}

func validateConfig(config Config) error {
	// Check if the config is valid and can be used for backups
	// TODO maybe skip missing sources via param
	log.Println("Validating config!")

	for _, unit := range config.units {
		if !unit.enabled {
			log.Printf("Unit '%s' is disabled. Skip validation for this unit!", unit.name)
			continue
		}

		for _, sourcePath := range unit.sources {
			// Each source path must be an existing directory!
			if !validatePath(sourcePath, true) {
				log.Printf("The given source path ('%s') does not exist or is no directory!", sourcePath)
				return errors.New("can't access source directory")
			}
		}
		// Also the destination path must exist!
		if !validatePath(unit.destination, true) {
			log.Printf("The given destination path ('%s') does not exist or is no directory!", unit.destination)
			return errors.New("can't access destination directory")
		}

		log.Printf("Unit '%s' is valid!", unit.name)
	}
	return nil
}

func readConfig(configPath string) (Config, error) {
	// Read config file at configPath
	log.Printf("Trying to read config file '%s'!", configPath)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Can't read config file! Exiting!")
		os.Exit(1)
	}

	// Read config file to Config struct
	c := Config{}
	c, err = c.FromYaml(data)

	if err != nil {
		return c, err
	}

	validateErr := validateConfig(c)
	return c, validateErr
}

func addFileToTar(tw *tar.Writer, path string, pathInArchive string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Format = tar.FormatGNU
		header.Name = pathInArchive
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zw *zip.Writer, path string, pathInArchive string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		header, headerErr := zip.FileInfoHeader(stat)
		if headerErr != nil {
			return headerErr
		}
		header.Method = zip.Deflate
		header.Name = pathInArchive
		// write the header to the zip archive
		writer, headerErr := zw.CreateHeader(header)
		if headerErr != nil {
			return err
		}
		// copy the file data to the zip
		if _, err := io.Copy(writer, file); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func main() {
	parser := argparse.NewParser("backmeup", "The lightweight backup tool for the CLI")
	parser.ExitOnHelp(true)
	configPath := parser.String("c", "config", &argparse.Options{Required: true, Help: "Path to the config.yml file", Default: "config.yml"})
	VERBOSE = *parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "Enable verbose logging", Default: false})
	DEBUG = *parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Enable debug logging", Default: false})

	if err := parser.Parse(os.Args); err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	conf, err := readConfig(*configPath)
	if err != nil {
		log.Println("Error while parsing yaml config!")
		os.Exit(1)
	}

	log.Println("Starting backup...")
	runBackup(conf)
}
