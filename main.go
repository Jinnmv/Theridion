// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {

	configuration := Configuration{}
	err := configuration.LoadConfigFromFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println("Config: ", *configuration)

	feedConfigs, err := InitFeedsConfiguration(configuration.Feeds.Path)
	if err != nil {
		log.Fatalln("error when reading feed configuration: ", err)
	}

	//price := Price{}
	//err = price.Fill(feedConfigs)

	//asyncHttpGet(feedConfigs)

	//Подготовим каналы и балансировщик
	feeds := make(chan *FeedData)
	quit := make(chan bool)
	b := new(Balancer)
	b.init(feeds)

	//Приготовимся перехватывать сигнал останова в канал keys
	keys := make(chan os.Signal, 1)
	signal.Notify(keys, os.Interrupt)

	//Запускаем балансировщик и генератор
	go b.balance(quit)
	go asyncHttpGet(feeds, feedConfigs)

	log.Printf("Started!")

	//Основной цикл программы:
	for {
		select {
		case <-keys: //пришла информация от нотификатора сигналов:
			log.Println("CTRL-C: Ожидаю завершения активных загрузок")
			quit <- true //посылаем сигнал останова балансировщику

		case <-quit: //пришло подтверждение о завершении от балансировщика
			log.Println("Загрузки завершены!")
			return
		}
	}

	//log.Println(feedConfig[0].Url)

}
