package main

import (
	"log"
	"regexp"
	_ "strconv"
	"sync"
)

type Worker struct {
	feeds   chan *FeedConfig // канал для заданий
	pending int              // кол-во оставшихся задач
	index   int              // позиция в куче
	id      int
	wg      *sync.WaitGroup //указатель на группу ожидания
}

type Mapping map[string]interface{}

func (w *Worker) work(done chan *Worker) {
	for {
		feed := <-w.feeds //читаем следующее задание
		w.wg.Add(1)       //инкриминируем счетчик группы ожидания
		parse(feed)       //загружаем файл
		w.wg.Done()       //сигнализируем группе ожидания что закончили
		done <- w         //показываем что завершили работу
	}
}

//Загрузка изображения
func parse(feed *FeedConfig) {
	//products := []*Price{}

	rg := *regexp.MustCompile(feed.Regex)

	match := rg.FindAllStringSubmatch(string(feed.Html), -1)

	if match == nil {
		log.Printf("[DEBUG]: PARSER nothing matched on", feed.Url)
		return
	}

	/*for _, goods := range match {

		price := Price{}
		for i, name := range rg.SubexpNames() {

			// Ignore the whole regexp match and unnamed groups
			if i == 0 || name == "" {
				continue
			}

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
				price.InStock = 0
			}
		}

		price.MarketName = feed.MarketName
		price.Category = feed.Define["category"]
		price.SubCategory = feed.Define["SubCategory"]

		products = append(products, &price)
	}

	log.Printf("[DEBUG]: PARSER %+v", len(products))*/
	log.Printf("[DEBUG]: PARSER %s, count: %d", feed.Url, len(match))
	return
}

func Deriviate(mapping Mapping, key string) interface{} {
	value, ok := mapping[key]
	if !ok {
		return key
	}
	return value
}
