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
	InStock      string
	URL          string
	ImageURL     string
	UpdateDate   time.Time
}

func NewPrice() *Price {
	return &Price{}
}

func (price *Price) Defaulting(defaultings map[string]string) {

	for key, value := range defaultings {
		switch key {
		case "name":
			price.Name = value
		case "category":
			price.Category = value
		case "subCategory":
			price.SubCategory = value
		case "manufacturer":
			price.Manufacturer = value
		case "scale":
			price.Scale = value
		case "currency":
			price.Currency = value
		case "sku":
			price.Sku = value
		case "inStock":
			price.InStock = value
		case "url":
			price.URL = value
		case "imgUrl":
			price.ImageURL = value
		}
	}
}

func (price *Price) Mapping(mappings map[string]map[string]string, data, keys []string) {

	for i, name := range keys {

		switch name {
		case "name":
			price.Name = Map(mappings[name], data[i])
		case "url":
			price.URL = Map(mappings[name], data[i])
		case "imgUrl":
			price.ImageURL = Map(mappings[name], data[i])
		case "manufacturer":
			price.Manufacturer = Map(mappings[name], data[i])
		case "sku":
			price.Sku = Map(mappings[name], data[i])
		case "price":
			x, err := strconv.ParseUint(data[i], 10, 0)
			if err == nil {
				price.Price = uint(x)
			}
		case "currency":
			price.Currency = Map(mappings[name], data[i])
		case "scale":
			price.Scale = Map(mappings[name], data[i])
		case "inStock":
			price.InStock = Map(mappings[name], data[i])
		}

	}

}

type PriceList []*Price

// Builder
func (products PriceList) Parse(feed *FeedConfig) PriceList {

	rg := *regexp.MustCompile(feed.Regex)

	match := rg.FindAllStringSubmatch(string(feed.Html), -1)

	if match == nil {
		log.Printf("[DEBUG]: PARSER nothing matched on", feed.Url)
		return nil
	}

	for _, goods := range match {

		price := NewPrice()
		price.MarketName = feed.MarketName

		price.Defaulting(feed.Defaulting)

		price.Mapping(feed.Mapping, goods[1:], rg.SubexpNames()[1:])

		price.UpdateDate = time.Now()

		/*for i, name := range rg.SubexpNames() {

			// Ignore the whole regexp match and unnamed groups
			if i == 0 || name == "" {
				continue
			}

			switch name {
			case "name":
				price.Name = Map(feed.Mapping[name], goods[i])
			case "url":
				price.URL = Map(feed.Mapping[name], goods[i])
			case "imgUrl":
				price.ImageURL = Map(feed.Mapping[name], goods[i])
			case "manufacturer":
				price.Manufacturer = Map(feed.Mapping[name], goods[i])
			case "sku":
				price.Sku = Map(feed.Mapping[name], goods[i])
			case "price":
				x, err := strconv.ParseUint(goods[i], 10, 0)
				if err == nil {
					price.Price = uint(x)
				}
			case "currency":
				price.Currency = Map(feed.Mapping[name], goods[i])
			case "scale":
				price.Scale = Map(feed.Mapping[name], goods[i])
			case "inStock":
				price.InStock = Map(feed.Mapping[name], goods[i])
			}
		}*/

		products = append(products, price)
	}

	log.Printf("[DEBUG]: PARSER count %+v", len(products))
	log.Printf("[DEBUG]: PARSER first: %+v", products[0])

	return products
}

func Map(mappings map[string]string, key string) string {

	value, ok := mappings[key]
	if !ok {
		return key
	}
	return value
}

// Find Dublicates TODO: implement
