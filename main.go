// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"github.com/Jinnmv/Theridion/configuration"
	"github.com/Jinnmv/Theridion/feedManager"
	"log"
)

func main() {

	configuration, err := configuration.New("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println("Config: ", *configuration)

	feedConfig, err := feedManager.New(configuration.Feeds.Path)
	if err != nil {
		log.Fatalln("error when reading feed configuration: ", err)
	}
	//fmt.Println(feedConfig)

	for i, feedConfItem := range *feedConfig {
		log.Printf("[%d]: %s\n", i, feedConfItem)
	}

}
