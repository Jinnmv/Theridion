package feedManager

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func (feedConfig *FeedConfig) FillConfig(fileName string) (*FeedConfig, error) {

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
		_, err = feedConfig.FillConfig(path.Join(feedsDir, f.Name()))
		if err != nil {
			log.Printf("Feed Config Manager: error when reading config file %s: %s", path.Join(feedsDir, f.Name()), err)
			//return nil, err
		}
	}

	return &feedConfig, nil

}
