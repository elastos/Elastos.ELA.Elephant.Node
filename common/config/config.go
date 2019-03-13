package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	DefaultConfigFilename = "./elephant.json"
)

var (
	Parameters ConfigFile
)

type Configuration struct {
	HttpRestPort int `json:"HttpRestPort"`
}

type ConfigFile struct {
	ConfigFile Configuration `json:"Configuration"`
}

func init() {
	if _, err := os.Stat(DefaultConfigFilename); os.IsNotExist(err) {
		log.Fatal("Can not find Elephant config file")
		os.Exit(1)
	} else {
		file, e := ioutil.ReadFile(DefaultConfigFilename)
		if e != nil {
			log.Fatalf("File error: %v\n", e)
			os.Exit(1)
		}
		// Remove the UTF-8 Byte Order Mark
		file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

		if e := json.Unmarshal(file, &Parameters); e != nil {
			log.Fatalf("Unmarshal json file erro %v", e)
			os.Exit(1)
		}
	}
}
