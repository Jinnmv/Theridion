// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database DbConnection `json:"database"`

	Feeds struct {
		Path string `json:"path"`
	} `json:"feeds"`

	Http struct {
		Threads byte `json:"threads"`
		Buffer  byte `json:"buffer"`
	} `json:"http"`

	Workers struct {
		Count    byte `json:"count"`
		Capacity byte `json:"capacity"`
	} `json:"workers"`
}

type DbConnection struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	DBName   string `json:"dbName"`
	Dialect  string `json:"dialect"`
	DSN      string `jsom:dsn`
}

var configInstance *Config

func GetConfigInstance(fileName string) (*Config, error) {

	if configInstance == nil {
		var err error
		configInstance, err = NewConfig(fileName)
		if err != nil {
			return nil, err
		}

	}
	return configInstance, nil
}

func NewConfig(fileName string) (*Config, error) {

	configInstance := Config{}

	if len(fileName) == 0 {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&configInstance)
	if err != nil {
		return nil, err
	}

	return &configInstance, nil
}
