package main

import (
	_ "errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func loader(ch chan *FeedConfig, feeds Feeds) {
	//chFeedConfig := make(chan *FeedConfig)
	wg := new(sync.WaitGroup)

	//feedData := []*FeedData{}

	// send all feedConfigs to channel
	/*go func (chFeedConfig *FeedConfig, feedConfigs []*FeedConfig){
		for {
			feedConfig, err := feedConfigs.Pop()
			if err != nil {
				return
			}
		}
	}(chFeedConfig, feedConfigs)*/

	for _, feedConfig := range feeds {
		wg.Add(1)

		// TODO: put anonymous func to separate
		go func(feedConfig *FeedConfig, wg *sync.WaitGroup) {
			defer wg.Done()
			defer timeTrack(time.Now(), feedConfig.Url)
			log.Printf("[INFO]: Fetching url [%s]", feedConfig.Url)

			resp, err := http.Get(feedConfig.Url)
			if err != nil {
				log.Printf("[ERROR]: Error when fetching url [%s]: %s\n", feedConfig.Url, err)
				return
			}
			defer resp.Body.Close()

			// Check for correct status code
			if resp.StatusCode < 200 || resp.StatusCode >= 400 {
				log.Printf("[ERROR]: [%s] Get request status code: %d\n", feedConfig.Url, resp.StatusCode)
				return
			}

			feedConfig.Html, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("[ERROR]: Error when fetching url source [%s]: %s\n", feedConfig.Url, err)
				return
			}

			ch <- feedConfig

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

	return

}
