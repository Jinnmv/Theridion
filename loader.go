package main

import (
	_ "errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func loader(ch chan *FeedConfig, feeds *Feeds) {
	feedConfigCh := make(chan *FeedConfig)
	wg := new(sync.WaitGroup)

	for i := 0; i < 2; i++ { // TODO: replace with value from config
		wg.Add(1)
		go downloader(feedConfigCh, ch, wg)

	}

	// Send all feedConfigs to channel
	for _, feedConfig := range *feeds {
		feedConfigCh <- feedConfig
	}

	close(feedConfigCh)
	wg.Wait()
	log.Println("[DEBUG]: LOADER All feeds are fetched")
	ch <- nil
	close(ch)

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

func downloader(feedConfigCh chan *FeedConfig, outCh chan *FeedConfig, wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		// Fetch next Feed Config from the channel
		feedConfig, ok := <-feedConfigCh
		if !ok {
			return
		}

		//log.Printf("[INFO]: Fetching url [%s]", feedConfig.Url)

		defer timeTrack(time.Now(), feedConfig.Url)
		resp, err := http.Get(feedConfig.Url)
		if err != nil {
			log.Printf("[ERROR]: Error when fetching url [%s]: %s\n", feedConfig.Url, err)
			continue
		}
		defer resp.Body.Close()

		// Check for correct status code
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			log.Printf("[ERROR]: [%s] Get request status code: %d\n", feedConfig.Url, resp.StatusCode)
			continue
		}

		feedConfig.Html, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR]: Error when fetching HTML from [%s]: %s\n", feedConfig.Url, err)
			continue
		}

		outCh <- feedConfig
	}
}
