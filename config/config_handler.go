package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/unsubble/tiny-url/common"
)

func ReadConfig(configPath string) *common.Config {
	file, err := os.Open(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("The config file was not found")
		}
		if errors.Is(err, os.ErrPermission) {
			log.Fatalf("The config file cannot be opened. Check permissions")
		}
		log.Fatalf("An error occurred while opening the config file: %v", err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	config := &common.Config{}
	err = decoder.Decode(config)

	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return config
}
