// Copyright 2014 Anton Savenko. All rights reserved.
// Use of this source code is governed by a GNU GPL
// license that can be found in the LICENSE file.

//

package main

import (
	"fmt"
	"github.com/Jinnmv/Theridion/configuration"
	"github.com/Jinnmv/Theridion/feedConfig"
	"log"
)

func main() {

	configuration, err := configuration.New("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Config: ", *configuration)

	feedConfig, err := feedConfig.New(configuration.Feeds.Path)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(feedConfig)
	fmt.Println(len(*feedConfig))

}
