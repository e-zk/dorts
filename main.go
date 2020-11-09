package main

import (
	"flag"
	"github.com/e-zk/dorts/templates"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"strings"
)

const defaultConfigDir = "{{ .xdgConfHome }}/dorts"
const defaultConfigPath = "{{ .confDir }}/dorts.toml"

var ignoredKeys = []string{"enabled", "path"}

var config *toml.Tree
var commonConfig map[string]string
var dorts []string
var configDir string
var verbose bool

func keyIsIgnored(key string) bool {
	for _, ignoredKey := range ignoredKeys {
		if key == ignoredKey {
			return true
		}
	}
	return false
}

/* substitute path with stuff */
func subsPath(tmpPath string) string {
	return strings.Replace(tmpPath, "~", os.Getenv("HOME"), 2)
}

func loadConfig() (*toml.Tree, error) {
	var err error
	var configPath string
	var conf *toml.Tree
	var vars = make(map[string]interface{})

	homeDir := os.Getenv("HOME")

	// set interface vars
	vars["xdgConfHome"] = os.Getenv("XDG_CONFIG_HOME")

	// if XDG_CONFIG_HOME is empty, use $HOME/.config instead
	if vars["xdgConfHome"].(string) == "" {
		vars["xdgConfHome"] = homeDir + "/.config"
	}

	// if the config dir is not already defined...
	if configDir == "" {
		// execute template on default config dir spec
		configDir, err = templates.ProcessString(defaultConfigDir, vars)
		if err != nil {
			return conf, err
		}

		// use environment variable for config dir if it exists
		dortsDirEnv := os.Getenv("DORTS_DIR")
		if len(dortsDirEnv) > 1 {
			configDir = os.Getenv("DORTS_DIR")
		}
	}

	// set confdir var
	vars["confDir"] = configDir

	// execute tempalte on default config file spec
	configPath, err = templates.ProcessString(defaultConfigPath, vars)
	if err != nil {
		return conf, err
	}

	return toml.LoadFile(configPath)
}

/* setup & parse command-line flags */
func parseFlags() {
	flag.StringVar(&configDir, "c", "", "path to config directory")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()
}

func main() {
	// setup config vars
	commonConfig = make(map[string]string)

	// logging
	log.SetFlags(log.Lmsgprefix | log.Lshortfile)

	parseFlags()

	// load config
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// todo check common
	if !config.Has("common") {
		log.Fatal("mandatory `common' table does not exist.")
	}

	// load keys/dorts
	for _, k := range config.Keys() {
		// load common table into commonConfig
		if k == "common" {
			common := config.Get("common").(*toml.Tree)
			for _, ck := range common.Keys() {
				commonConfig[ck] = common.Get(ck).(string)

			}
		} else {
			dorts = append(dorts, k)
		}
	}

	// for each dort read template:
	for _, dort := range dorts {
		var outputPath string
		vars := make(map[string]interface{})
		templatePath := configDir + "/" + dort + ".tmpl"

		// get dort-specific settings
		dortConfig := config.Get(dort).(*toml.Tree)

		// skip disabled dorts
		if dortConfig.Has("enabled") && dortConfig.Get("enabled").(bool) == false {
			log.Printf("dort `%s' is disabled. skipping.\n", dort)
			continue
		}

		// skip dorts without a template file
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			log.Printf("template for dort `%s' does not exist. skipping.\n", dort)
			continue
		}

		// check & parse path
		if !dortConfig.Has("path") {
			log.Fatalf("dort `%s' does not have mandatory `path' key.\n", dort)
		} else {
			originalPath := dortConfig.Get("path").(string)
			outputPath = subsPath(originalPath)
		}

		// construct interface to execute on template
		for k, v := range commonConfig {
			// add global variables
			vars[k] = v
		}

		// add 'local' keys to interface
		// this overrides global settings
		for _, k := range dortConfig.Keys() {
			if keyIsIgnored(k) {
				continue
			}
			vars[k] = dortConfig.Get(k).(string)
		}

		// parse tempalte
		result, err := templates.ProcessFile(templatePath, vars)
		if err != nil {
			log.Println("error parsing template.")
			log.Fatal(err)
		}

		// open file
		f, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 00644)
		if err != nil {
			log.Fatal(err)
		}

		// write to file
		f.WriteString(result)

		// close file
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
