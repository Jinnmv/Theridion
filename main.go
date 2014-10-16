// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"github.com/Jinnmv/Theridion/configuration"
	"github.com/Jinnmv/Theridion/feedManager"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	configuration, err := configuration.New("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println("Config: ", *configuration)

	feedConfigs, err := feedManager.New(configuration.Feeds.Path)
	if err != nil {
		log.Fatalln("error when reading feed configuration: ", err)
	}

	price := Price{}
	err = price.Fill(feedConfigs)

	//log.Println(feedConfig[0].Url)

}

func downloader(out chan string, url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	out <- string(body)
}
