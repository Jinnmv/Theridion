package feedConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type FeedConfig []struct {
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

func (feedConfig *FeedConfig) fillConfig(fileName string) (*FeedConfig, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	newFeedConfig := FeedConfig{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&newFeedConfig)
	if err != nil {
		return nil, err
	}

	fmt.Println(newFeedConfig)
	*feedConfig = append(*feedConfig, newFeedConfig...)
	return &newFeedConfig, nil
}

func New(feedsDir string) (*FeedConfig, error) {

	feedConfig := FeedConfig{}

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		_, err = feedConfig.fillConfig(path.Join(feedsDir, f.Name()))
		if err != nil {
			return nil, err
		}
		//fmt.Println(f.Name())
	}

	return &feedConfig, nil

}
