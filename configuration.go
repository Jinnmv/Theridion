// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
		DBName   string `json:"dbName"`
		Dialect  string `json:"dialect"`
	} `json:"database"`

	Feeds struct {
		Path string `json:"path"`
	} `json:"feeds"`

	Http struct {
		Threads byte `json:"threads"`
	} `json:"http"`

	Workers struct {
		Count    byte `json:"count"`
		Capacity byte `json:"capacity"`
	} `json:"workers"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (configuration *Configuration) LoadFromFile(fileName string) error {

	if len(fileName) == 0 {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	//configuration := *Configuration{}

	err = decoder.Decode(configuration)
	if err != nil {
		return err
	}

	return nil
}
