package main

import (
	"github.com/Jinnmv/Theridion/feedManager"
	"log"
	"time"
)

type Price []struct {
	Name         string
	Category     string
	SubCategory  string
	Manufacturer string
	Scale        string
	Price        uint
	Currency     string
	Sku          string
	Market       string
	InStock      byte
	URL          string
	ImageURL     string
	UpdateDate   time.Time
}

func (price *Price) Fill(feedConfigs *feedManager.FeedConfigs) error {
	for i, feedConfig := range *feedConfigs {
		log.Printf("[%d]: %s\n", i, feedConfig)
	}
	return nil
}
