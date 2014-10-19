package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type FeedConfig struct {
	MarketName string            `json:"market`
	Url        string            `json:"url"`
	DataFields map[string]string `json:"dataFields`
	Parse      struct {
		Regex           string            `json:"regex"`
		RegexDataFields map[string]string `json:"regexDataFields"`
	} `json:"parse"`
	Derivations struct {
		Mapping map[string]map[string]interface{} `json:"mapping"`
	} `json:"derivations"`
}

func readFeedConfig(feedConfigs []*FeedConfig, fileName string) ([]*FeedConfig, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	newFeedConfigs := []*FeedConfig{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&newFeedConfigs)
	if err != nil {
		return nil, err
	}

	//feedConfigs = append(feedConfigs, newFeedConfigs...)
	return newFeedConfigs, nil
}

func InitFeedsConfiguration(feedsDir string) ([]*FeedConfig, error) {

	feedConfigs := []*FeedConfig{}

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		res, err := readFeedConfig(feedConfigs, path.Join(feedsDir, file.Name()))
		if err != nil {
			log.Printf("[ERROR]: Feed Config Manager: error when reading config file %s: %s", path.Join(feedsDir, file.Name()), err)
		} else {
			feedConfigs = append(feedConfigs, res...)
		}

	}

	return feedConfigs, nil
}
