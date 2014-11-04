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

	config := NewConfiguration()
	err := config.LoadFromFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	feeds := NewFeeds()
	err = feeds.LoadFromDir(config.Feeds.Path)
	if err != nil {
		log.Fatalln("Error when reading feed configuration: ", err)
	}

	log.Println("[DEBUG]: Feeds count", len(feeds))

	//price := Price{}
	//err = price.Fill(feedConfigs)

	//Init Channels and Balancer
	feedsCh := make(chan *FeedConfig, config.Http.Buffer)
	quitCh := make(chan bool)
	balancer := Balancer{}
	balancer.Init(feedsCh, config.Workers.Count, config.Workers.Capacity)

	//Init OS signal interceptor ot channel keys
	keys := make(chan os.Signal, 1)
	signal.Notify(keys, os.Interrupt)

	//Run Balancer and Loader
	go balancer.Balance(quitCh)

	downloader := NewDownloader()
	downloader.Init(feeds, config.Http.Threads)
	go downloader.Load(feedsCh)

	log.Printf("Started!")

	//Main cycle:
	for {
		select {
		case <-keys: //пришла информация от нотификатора сигналов:
			log.Println("CTRL-C: Ожидаю завершения активных загрузок")
			quitCh <- true //посылаем сигнал останова балансировщику

		case <-quitCh: //пришло подтверждение о завершении от балансировщика
			log.Println("Загрузки завершены!")
			return
		}
	}

	//log.Println(feedConfig[0].Url)

}
