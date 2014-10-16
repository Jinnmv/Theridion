package feedManager

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

type FeedConfigs []FeedConfig // TODO: use of type FeedConfigs []*FeedConfig

func (feedConfigs *FeedConfigs) FillConfig(fileName string) (*FeedConfigs, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	newFeedConfigs := FeedConfigs{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&newFeedConfigs)
	if err != nil {
		return nil, err
	}

	*feedConfigs = append(*feedConfigs, newFeedConfigs...)
	return &newFeedConfigs, nil
}

func New(feedsDir string) (*FeedConfigs, error) {

	feedConfigs := FeedConfigs{}

	files, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		_, err = feedConfigs.FillConfig(path.Join(feedsDir, f.Name()))
		if err != nil {
			log.Printf("Feed Config Manager: error when reading config file %s: %s", path.Join(feedsDir, f.Name()), err)
			//return nil, err
		}
	}

	return &feedConfigs, nil

}
