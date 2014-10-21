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

	Parse struct {
		Regex           string            `json:"regex"`
		RegexDataFields map[string]string `json:"regexDataFields"`
	} `json:"parse"`

	Derivations struct {
		Mapping map[string]map[string]interface{} `json:"mapping"`
	} `json:"derivations"`

	Html []byte `json:"-"`
}

type Feeds []*FeedConfig

func NewFeeds() Feeds {
	return Feeds{}
}

func (feeds *Feeds) LoadFromFile(fileName string) error {

	newFeeds := []FeedConfig{}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&newFeeds)
	if err != nil {
		return err
	}

	for _, feed := range newFeeds {
		*feeds = append(*feeds, &feed)
	}

	return nil
}

func (feeds *Feeds) LoadFromDir(feedsDir string) error {

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := feeds.LoadFromFile(path.Join(feedsDir, file.Name()))
		if err != nil {
			log.Printf("[ERROR]: Feed Config Manager: error when reading config file %s: %s", path.Join(feedsDir, file.Name()), err)
		}

	}

	return nil
}

// define next mehods:
// * NewFeedConfiguration - constructor
// * (*FeedConfigList) LoadFromFile
// * (*FeedConfigList) LoadFromDir
// * (*FeedConfigList) Pop
// * (*FeedConfigList) Push
// * (*FeedConfigList) Append
