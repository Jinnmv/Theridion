package main

import (
	_ "errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type FeedData struct {
	FeedConfig *FeedConfig	// TODO: *FeedConfig
	Html       []byte
}

func loader(ch chan *FeedData, feedConfigs []*FeedConfig) []*FeedData {
	//ch := make(chan *FeedData)
	wg := new(sync.WaitGroup)

	feedData := []*FeedData{}
	
	// for { send all feedConfigs to channel }

	for _, feedConfig := range feedConfigs {
		wg.Add(1)
		go func(feedConfig *FeedConfig, wg *sync.WaitGroup) {
			defer wg.Done()
			log.Printf("[INFO]: Fetching url [%s]", feedConfig.Url)

			resp, err := http.Get(feedConfig.Url)
			if err != nil {
				log.Printf("[ERROR]: Error when fetching url [%s]: %s\n", feedConfig.Url, err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("[ERROR]: Error when fetching url source [%s]: %s\n", feedConfig.Url, err)
				return
			}

			ch <- &FeedData{feedConfig, body}

		}(feedConfig, wg)

	}

	go func() {
		wg.Wait()
		log.Println("[DEBUG]: LOADER All feeds are fetched")
		ch <- nil
		close(ch)
	}()

	/*for {
		select {
		case r, ok := <-ch:
			if ok {
				log.Printf("[INFO]: %s was fetched %d bytes", r.FeedConfig.Url, len(r.Html))
				feedData = append(feedData, r)
			} else {
				ch = nil
				log.Println("[DEBUG]: chanel is not ok")
			}

			if len(feedData) == len(feedConfigs) {
				return feedData
			}
		}

	}*/

	return feedData

}
