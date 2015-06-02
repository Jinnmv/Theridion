package main

import (
	"database/sql"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Price struct {
	ID           uint64         `db:"id"`
	Name         string         `db:"name"`
	Category     sql.NullString `db:"category"`
	SubCategory  sql.NullString `db:"sub_category"`
	Manufacturer string         `db:"manufacturer"`
	Scale        string         `db:"scale"`
	Price        int            `db:"price"` //uint?
	Currency     string         `db:"currency"`
	Sku          string         `db:"sku"`
	MarketName   string         `db:"market_name"`
	InStock      string         `db:"in_stock"`
	URL          string         `db:"url"`
	ImageURL     string         `db:"image_url"`
	UpdateDate   time.Time      `db:"update_date"`
}

func NewPrice() *Price {
	return &Price{}
}

func (price *Price) Defaulting(defaultings map[string]string) {

	for key, value := range defaultings {
		valid := len(value) != 0
		switch key {
		case "name":
			price.Name = value
		case "category":
			price.Category = sql.NullString{String: value, Valid: valid}
		case "subCategory":
			price.SubCategory = sql.NullString{String: value, Valid: valid}
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
			x, err := strconv.Atoi(data[i])
			if err == nil {
				price.Price = x
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

func (price *Price) TrimName() { // TODO: implement sending a field required to trim
	price.Name = strings.TrimSpace(price.Name)
}

func (p *Price) CleanUrl() {
	p.URL = strings.Split(p.URL, "?")[0]
}

func (p *Price) NormalizeManufacturer() {
	p.Manufacturer = strings.Title(strings.ToLower(p.Manufacturer))
}

func (p *Price) EnrichImageURL() {

	p.ImageURL = strings.TrimPrefix(p.ImageURL, "/")

	if !strings.HasPrefix(p.ImageURL, "http") {
		u, err := url.Parse(p.URL)
		if err != nil {
			log.Println("[WARNING]: Unable to parse Host URL:", p.URL)
			return
		}
		p.ImageURL = strings.Join([]string{u.Scheme, "://", u.Host, "/", p.ImageURL}, "")
	}

}

func (p *Price) CleanEmptySubCategory() {

}

type PriceCollection []*Price

func NewPriceCollection() *PriceCollection {
	return &PriceCollection{}
}

// Builder
func (pc *PriceCollection) Parse(feed *Feed) *PriceCollection {

	//defer timeTrack(time.Now(), "[TIMER] parsing")

	rg := *regexp.MustCompile(feed.Regex)

	match := rg.FindAllStringSubmatch(string(feed.Html), -1)

	if match == nil {
		log.Printf("[DEBUG]: PARSER nothing matched on %s", feed.Url)
		return nil
	}

	for _, goods := range match {

		price := NewPrice()
		price.MarketName = feed.MarketName

		price.Defaulting(feed.Defaulting)

		price.Mapping(feed.Mapping, goods[1:], rg.SubexpNames()[1:]) // skipping 1-st element as it's a whole string

		price.CleanUrl()
		price.EnrichImageURL()
		price.TrimName()
		price.NormalizeManufacturer()

		price.UpdateDate = time.Now()

		*pc = append(*pc, price)
	}

	log.Printf("[DEBUG]: PARSER count %+v", len(*pc))

	return pc
}

func Map(mappings map[string]string, key string) string {

	value, ok := mappings[key]
	if !ok {
		return key
	}
	return value
}

// Find Dublicates TODO: implement

type PriceManager struct {
	storage   Storage
	PriceList PriceCollection
}

func NewPriceManager() *PriceManager {
	pm := PriceManager{}

	return &pm
}

func (pm *PriceManager) SetStorage(storage Storage) {
	pm.storage = storage
}

func (pm *PriceManager) Write() (int, error) {

	pl := make([]*interface{}, len(pm.PriceList)) // create slice container of interface{}
	for i, v := range pl {                        // fill in content
		pl[i] = v
	}

	return pm.storage.Write(pl)
}

func (pm *PriceManager) Load() error {

	return nil //TODO: implement
}
