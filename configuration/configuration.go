// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package configuration

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
		Workers byte `json:"workers"`
	} `json:"http"`
}

func New(fileName string) (*Configuration, error) {

	if len(fileName) == 0 {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configuration := Configuration{}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
