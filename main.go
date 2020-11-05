package main

import (
	"bytes"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"strings"
	"text/template"
)

var ignoredKeys = []string{"enabled", "path"}

var config *toml.Tree
var commonConfig map[string]string
var dorts []string

var path = "/home/zzz/etc/dorts/dorts.toml"
var confDir = "/home/zzz/etc/dorts"

func keyIsIgnored(key string) bool {
	for _, ignoredKey := range ignoredKeys {
		if key == ignoredKey {
			return true
		}
	}
	return false
}

/* process template string */
func process(t *template.Template, vars interface{}) (string, error) {
	var err error = nil
	var tmpBytes bytes.Buffer

	err = t.Execute(&tmpBytes, vars)
	if err != nil {
		return "", err
	}

	return tmpBytes.String(), err
}

/* process template file */
func processFile(path string, vars interface{}) (string, error) {
	var err error = nil

	tmp, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	return process(tmp, vars)
}

/* substitute path with stuff */
func subsPath(path string) string {
	return strings.Replace(path, "~", os.Getenv("HOME"), 1)
}

func main() {
	// setup config vars
	commonConfig = make(map[string]string)

	// logging
	log.SetFlags(log.Lmsgprefix | log.Lshortfile)

	// load config
	config, err := toml.LoadFile(path)
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
		templatePath := confDir + "/" + dort + ".tmpl"

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
		result, err := processFile(templatePath, vars)
		if err != nil {
			log.Println("error parsing template.")
			log.Fatal(err)
		}

		// open file
		f, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 00644)
		if err != nil {
			log.Fatal(err)
		}

		f.WriteString(result)

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
