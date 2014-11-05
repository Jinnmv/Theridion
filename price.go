package main

import (
	"log"
	"regexp"
	"strconv"
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

/*func (price *Price) Fill(feedConfigs []*FeedConfig) error {
	for i, feedConfig := range feedConfigs {
		log.Printf("[%d]: %s\n", i, feedConfig)
	}
	return nil
}*/

type PriceList []*Price

func (products PriceList) Parse(feed *FeedConfig) PriceList {

	rg := *regexp.MustCompile(feed.Regex)

	match := rg.FindAllStringSubmatch(string(feed.Html), -1)

	if match == nil {
		log.Printf("[DEBUG]: PARSER nothing matched on", feed.Url)
		return nil
	}

	for _, goods := range match {

		price := Price{}
		for i, name := range rg.SubexpNames() {

			// Ignore the whole regexp match and unnamed groups
			if i == 0 || name == "" {
				continue
			}

			goods[i] = Derivate(feed.Derivations.Mapping[name], goods[i])

			switch name {
			case "name":
				price.Name = goods[i]
			case "url":
				price.URL = goods[i]
			case "imgUrl":
				price.ImageURL = goods[i]
			case "manufacturer":
				price.Manufacturer = goods[i]
			case "sku":
				price.Sku = goods[i]
			case "price":
				z, err := strconv.ParseUint(goods[i], 10, 64)
				if err == nil {
					price.Price = uint(z)
				}
			case "currency":
				price.Currency = goods[i]
			case "scale":
				price.Scale = goods[i]
			case "inStock":
				z, err := strconv.ParseUint(goods[i], 10, 8)
				if err == nil {
					price.InStock = byte(z)
				}
			}
		}

		price.MarketName = feed.MarketName
		price.Category = feed.Define["category"]
		price.SubCategory = feed.Define["SubCategory"]
		price.UpdateDate = time.Now()

		products = append(products, &price)
	}

	log.Printf("[DEBUG]: PARSER count %+v", len(products))
	log.Printf("[DEBUG]: PARSER first: %+v", products[0])

	return products
}

func Derivate(mappings map[string]string, key string) string {

	value, ok := mappings[key]
	if !ok {
		return key
	}
	return value
}

// Find Dublicates TODO: implement
