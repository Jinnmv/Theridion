package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type FeedConfig struct {
	MarketName string `json:"market"`
	Url        string `json:"url"`

	Regex string `json:"regex"`

	Defaulting map[string]string            `json:"defaulting"`
	Mapping    map[string]map[string]string `json:"mapping"`

	Html []byte `json:"-"`
}

type Feeds []*FeedConfig

var feedsInst *Feeds

func GetFeedsInstance() *Feeds {
	if feedsInst == nil {
		feedsInst = NewFeeds()
	}
	return feedsInst
}

func NewFeeds() *Feeds {
	feedsInst = &Feeds{}
	return feedsInst
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

/*func (f *Feeds) IndexOf(fc *FeedConfig) int {
	for i, item := range *f {
		if item == fc {
			return i
		}
	}
	return -1
}

*/

func (f *Feeds) Pop() (*FeedConfig, error) {

	if len(*f) == 0 {
		return nil, errors.New("Can't pop empty stack")
	}

	fd := *f

	x := fd[len(*f)-1]
	*f = fd[:len(*f)-1]

	return x, nil
}

// define next mehods:
// * NewFeedConfiguration - constructor
// * (*FeedConfigList) LoadFromFile
// * (*FeedConfigList) LoadFromDir
// * (*FeedConfigList) Pop
// * (*FeedConfigList) Push
// * (*FeedConfigList) Append
