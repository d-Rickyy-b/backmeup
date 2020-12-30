package main

import (
	"backmeup/archiver"
	"backmeup/config"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/bmatcuk/doublestar/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var VERBOSE bool
var DEBUG bool

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

func handleExcludes(filePath string, excludePatterns []string) bool {
	// Checks if the path is excluded by any of the given exclude patterns
	for _, excludePattern := range excludePatterns {
		matched := handleExclude(filePath, excludePattern)

		if matched {
			return true
		}
	}

	return false
}

func getFiles(sourcePath string, excludes []string) ([]archiver.BackupFileMetadata, error) {
	// Returns all file paths recursively within a certain source directory
	var pathsToBackup []archiver.BackupFileMetadata

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

			fileMetadata := archiver.BackupFileMetadata{
				Path:           path,
				BackupBasePath: sourcePath,
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

func writeBackup(filesToBackup []archiver.BackupFileMetadata, unit config.Unit) {
	now := time.Now()
	timeStamp := now.Format("2006-01-02_15-04")
	backupBasePath := unit.Destination

	if unit.AddSubfolder {
		newBackupBasePath := filepath.Join(unit.Destination, unit.Name)
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

	backupArchiveName := unit.Name + "-" + timeStamp + "." + unit.ArchiveType
	backupArchivePath := filepath.Join(backupBasePath, backupArchiveName)

	// TODO check if archive already exists. If yes, append -1 to it and try again

	archiver.WriteArchive(backupArchivePath, filesToBackup, unit)
	log.Printf("Archive created successfully at '%s'", backupArchivePath)
}

func backupUnit(unit config.Unit) {
	// Start backup for a single unit. Each backup creates a single archive file
	if !unit.Enabled {
		log.Printf("Skipping backup for unit '%s' because it's disabled.\n", unit.Name)

		return
	}

	log.Printf("Creating backup for unit '%s'\n", unit.Name)

	var (
		filesToBackup    []archiver.BackupFileMetadata
		processedSources []string
	)

	// Check all source files from the disk in the specified source directories
	for _, sourcePath := range unit.Sources {
		sourcePath = filepath.Clean(sourcePath)

		// Prevent duplicate source paths
		for _, processedPath := range processedSources {
			if sourcePath == processedPath {
				log.Printf("Found duplicate source path '%s'. Skipping!", sourcePath)

				continue
			}
		}

		processedSources = append(processedSources, sourcePath)

		files, err := getFiles(sourcePath, unit.Excludes)
		if err != nil {
			log.Printf("Error for unit '%s' while reading directory '%s'! Skipping!", unit.Name, sourcePath)

			continue
		}

		filesToBackup = append(filesToBackup, files...)
	}

	if len(filesToBackup) == 0 {
		log.Printf("No files found for sources in unit '%s'. Creating no backup!", unit.Name)

		return
	}

	writeBackup(filesToBackup, unit)
}

func runBackup(config config.Config) {
	// Start the backup(s) defined in the config object
	for _, unit := range config.Units {
		backupUnit(unit)
	}
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

	conf, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Println("Error while parsing yaml config!")
		os.Exit(1)
	}

	log.Println("Starting backup...")
	runBackup(conf)
}
