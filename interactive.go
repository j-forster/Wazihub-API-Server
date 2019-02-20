package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func interactive() {

	if len(os.Args) == 4 {
		if os.Args[1] == "set" {
			err := setConfig(os.Args[2], os.Args[3])
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	if len(os.Args) == 3 {
		if os.Args[1] == "get" {
			str, err := getConfig(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			log.Println(str)
			return
		}
	}

	log.Fatal("Unknown arguments.")
}

var config map[string]string

func readConfig() error {
	data, err := ioutil.ReadFile("/etc/wazihub.json")
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		config = make(map[string]string)
		return nil
	} else {
		return json.Unmarshal(data, &config)
	}
}

func writeConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("/etc/wazihub.json", data, 0644)
}

func setConfig(name, value string) error {
	if config == nil {
		err := readConfig()
		if err != nil {
			return err
		}
	}
	config[name] = value
	return writeConfig()
}

func getConfig(name string) (string, error) {
	if config == nil {
		err := readConfig()
		if err != nil {
			return "", err
		}
	}
	return config[name], nil
}
