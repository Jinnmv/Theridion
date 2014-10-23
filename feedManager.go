package main

import (
	"encoding/json"
	_ "errors"
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

func (feeds *Feeds) AppendFromFile(fileName string) error {

	newFeeds := Feeds{}

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

	*feeds = append(*feeds, newFeeds...)

	return nil
}

func (feeds *Feeds) LoadFromDir(feedsDir string) error {

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := feeds.AppendFromFile(path.Join(feedsDir, file.Name()))
		if err != nil {
			log.Printf("[ERROR]: Feed Config Manager: error when reading config file %s: %s", path.Join(feedsDir, file.Name()), err)
		}

	}

	return nil
}

/*func (feeds *Feeds) Pop() (*FeedConfig, error) {
	if len(*feeds) == 0 {
		return nil, errors.New("Can't pop empty stack")
	}

	fd := *feeds

	x := fd[len(*feeds)-1]
	*feeds = fd[:len(*feeds)-1]

	return x, nil
}*/

// define next mehods:
// * NewFeedConfiguration - constructor
// * (*FeedConfigList) LoadFromFile
// * (*FeedConfigList) LoadFromDir
// * (*FeedConfigList) Pop
// * (*FeedConfigList) Push
// * (*FeedConfigList) Append
