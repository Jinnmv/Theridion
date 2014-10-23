package main

import (
	"log"
	"time"
)

type Price struct {
	Name         string
	Category     string
	SubCategory  string
	Manufacturer string
	Scale        string
	Price        uint
	Currency     string
	Sku          string
	MarketName   string
	InStock      byte
	URL          string
	ImageURL     string
	UpdateDate   time.Time
}

func (price *Price) Fill(feedConfigs []*FeedConfig) error {
	for i, feedConfig := range feedConfigs {
		log.Printf("[%d]: %s\n", i, feedConfig)
	}
	return nil
}

// TODO: temporary
func ParsePrice(feedConfig *FeedConfig) {

}
