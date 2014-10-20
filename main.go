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

	configuration, err := NewConfiguration("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	feedConfigs, err := InitFeedsConfiguration(configuration.Feeds.Path)
	if err != nil {
		log.Fatalln("error when reading feed configuration: ", err)
	}

	//price := Price{}
	//err = price.Fill(feedConfigs)

	//Init Channels and Balancer
	feeds := make(chan *FeedData)
	quit := make(chan bool)
	b := new(Balancer)
	b.init(feeds)

	//Init OS signal interceptor ot channel keys
	keys := make(chan os.Signal, 1)
	signal.Notify(keys, os.Interrupt)

	//Run Balancer and Loader
	go b.balance(quit)
	go loader(feeds, feedConfigs)

	log.Printf("Started!")

	//Main cycle:
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
