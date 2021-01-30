package config

import (
	"backmeup/bkperrors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Unit struct {
	Name             string
	Sources          []string
	Destination      string
	Excludes         []string
	ArchiveType      string
	AddSubfolder     bool
	Enabled          bool
	UseAbsolutePaths bool
}

type Config struct {
	Units []Unit
}

// Helper struct for parsing the yaml
type yamlUnit struct {
	Sources          *[]string `yaml:"sources"`
	Destination      *string   `yaml:"destination"`
	Excludes         *[]string `yaml:"excludes"`
	ArchiveType      *string   `yaml:"archive_type"`
	AddSubfolder     *bool     `yaml:"add_subfolder"`
	Enabled          *bool     `yaml:"enabled"`
	UseAbsolutePaths *bool     `yaml:"use_absolute_paths"`
}

func (config Config) FromYaml(yamlData []byte) (Config, error) {
	// Create a config object from yaml byte array
	unitMap := make(map[string]yamlUnit)

	log.Println("Parsing config yaml")

	unmarshalErr := yaml.Unmarshal(yamlData, &unitMap)
	if unmarshalErr != nil {
		log.Fatalf("Unmarshal error: %v", unmarshalErr)
	}

	// After parsing the yaml into unitMap, we iterate over all available units
	for unitName, yamlUnit := range unitMap {
		unit := Unit{}

		// Set defaults
		unit.Enabled = true
		if yamlUnit.Enabled != nil {
			unit.Enabled = *yamlUnit.Enabled
		}

		unit.AddSubfolder = false
		if yamlUnit.AddSubfolder != nil {
			unit.AddSubfolder = *yamlUnit.AddSubfolder
		}

		unit.ArchiveType = "tar.gz"
		if yamlUnit.ArchiveType != nil {
			unit.ArchiveType = *yamlUnit.ArchiveType
		}

		unit.Excludes = []string{}
		if yamlUnit.Excludes != nil {
			unit.Excludes = *yamlUnit.Excludes
		}

		unit.UseAbsolutePaths = true
		if yamlUnit.UseAbsolutePaths != nil {
			unit.UseAbsolutePaths = *yamlUnit.UseAbsolutePaths
		}

		if yamlUnit.Sources == nil || yamlUnit.Destination == nil {
			log.Fatalf("Sources or destination can't be parsed for unit '%s'", unitName)
		} else {
			unit.Sources = *yamlUnit.Sources
			unit.Destination = *yamlUnit.Destination
		}

		unit.Name = unitName

		config.Units = append(config.Units, unit)
	}

	return config, nil
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

func (config *Config) validate() error {
	// Check if the config is valid and can be used for backups
	// TODO maybe skip missing sources via param
	log.Println("Validating config!")

	for _, unit := range config.Units {
		if !unit.Enabled {
			log.Printf("Unit '%s' is disabled. Skip validation for this unit!", unit.Name)

			continue
		}

		for _, sourcePath := range unit.Sources {
			// Each source path must be an existing directory!
			if !validatePath(sourcePath, true) {
				log.Printf("The given source path ('%s') does not exist or is no directory!", sourcePath)

				return bkperrors.ErrCannotAccessSrcDir
			}
		}
		// Also the destination path must exist!
		if !validatePath(unit.Destination, true) {
			log.Printf("The given destination path ('%s') does not exist or is no directory!", unit.Destination)

			return bkperrors.ErrCannotAccessDstDir
		}

		log.Printf("Unit '%s' is valid!", unit.Name)
	}

	return nil
}

func ReadConfig(configPath string) (Config, error) {
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

	validateErr := c.validate()

	return c, validateErr
}
