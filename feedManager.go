package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Feed struct {
	MarketName string `json:"market"`
	Url        string `json:"url"`

	Regex string `json:"regex"`

	Defaulting map[string]string            `json:"defaulting"`
	Mapping    map[string]map[string]string `json:"mapping"`

	Html []byte `json:"-"`
}

type FeedCollection []*Feed

var feedColInst *FeedCollection

func GetFeedColInstance() *FeedCollection {
	if feedColInst == nil {
		feedColInst = NewFeedCollection()
	}
	return feedColInst
}

func NewFeedCollection() *FeedCollection {
	feedColInst = &FeedCollection{}
	return feedColInst
}

func (fc *FeedCollection) AppendFromFile(fileName string) error {

	newFeeds := FeedCollection{}

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

	*fc = append(*fc, newFeeds...)

	return nil
}

func (fc *FeedCollection) LoadFromDir(feedsDir string) error {

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := fc.AppendFromFile(path.Join(feedsDir, file.Name()))
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

func (fc *FeedCollection) Pop() (*Feed, error) {

	if len(*fc) == 0 {
		return nil, errors.New("Can't pop empty stack")
	}

	fd := *fc

	x := fd[len(*fc)-1]
	*fc = fd[:len(*fc)-1]

	return x, nil
}

// define next mehods:
// * NewFeedConfiguration - constructor
// * (*FeedConfigList) LoadFromFile
// * (*FeedConfigList) LoadFromDir
// * (*FeedConfigList) Pop
// * (*FeedConfigList) Push
// * (*FeedConfigList) Append
