// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"github.com/Jinnmv/Theridion/Balancer"
	"log"
	"os"
	"os/signal"
	"time"
)

var configFileName = "config.json"

func main() {

	defer timeTrack(time.Now(), "Full run took:")

	config, err := NewConfig(configFileName)
	if err != nil {
		log.Fatalln(err)
	}

	feeds := NewFeeds()
	err = feeds.LoadFromDir(config.Feeds.Path)
	if err != nil {
		log.Fatalln("Error when reading feed configuration: ", err)
	}

	log.Println("[DEBUG]: Feeds count", len(*feeds))

	// Init Storage
	storage := NewStorage(
		config.Database.Dialect,
		config.Database.Hostname,
		config.Database.Port,
		config.Database.DBName,
		config.Database.Username,
		config.Database.Password)
	defer storage.Close()

	//Init Channels and Balancer
	feedsChannel := make(chan interface{}, config.Http.Buffer)
	quitCh := make(chan bool)

	balancer := Balancer.NewBalancer()
	balancer.Init(feedsChannel, workerJob, config.Workers.Count, config.Workers.Capacity)

	//Init OS signal interceptor ot channel keys
	keys := make(chan os.Signal, 1)
	signal.Notify(keys, os.Interrupt)

	//Run Balancer and Loader
	go balancer.Balance(quitCh)

	downloader := NewDownloader(feeds, config.Http.Threads)
	go downloader.Load(feedsChannel)

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

func workerJob(feed interface{}) {
	priceList := NewPriceList()
	priceList.Parse(feed.(*FeedConfig))

	config, err := GetConfigInstance(configFileName)
	if err != nil {
		log.Fatalf("[FATAL]: Unable to load configuration: %v", err)
	}

	//Init DB
	storage := GetStorageInstance(
		config.Database.Dialect,
		config.Database.Hostname,
		config.Database.Port,
		config.Database.DBName,
		config.Database.Username,
		config.Database.Password)

	_, err = storage.Write(*priceList)
	if err != nil {
		log.Println("[ERROR]: DB error when inserting data", err)
	}

	// reduce memory

	//log.Println("[DEBUG]: JOB: Feed index:", GetFeedsInstance().IndexOf(feed.(*FeedConfig)))
	//feed(*FeedConfig)).Html = nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Printf("%s: %+v", msg, err)
	}
}
