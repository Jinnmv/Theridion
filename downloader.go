package main

import (
	_ "errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Downloader struct {
	feeds        *FeedCollection
	threads      byte
	feedConfigCh chan *Feed
	wg           *sync.WaitGroup
}

func NewDownloader(feeds *FeedCollection, threads byte) *Downloader {
	return &Downloader{feeds: feeds,
		threads:      threads,
		feedConfigCh: make(chan *Feed),
		wg:           new(sync.WaitGroup)}
}

func (d *Downloader) Load(outCh chan interface{}) {

	for i := 0; i < int(d.threads); i++ {
		d.wg.Add(1)
		go d.fetch(outCh)

	}

	// Send all feedConfigs to channel
	for {
		feedConfig, err := d.feeds.Pop()
		if err != nil {
			break
		}
		d.feedConfigCh <- feedConfig
	}

	close(d.feedConfigCh)
	d.wg.Wait()
	log.Println("[DEBUG]: LOADER All feeds are fetched")
	outCh <- nil
	close(outCh)

	return

}

func (d *Downloader) fetch(outCh chan interface{}) {

	defer d.wg.Done()

	for {

		// Fetch next Feed Config from the channel
		feedConfig, ok := <-d.feedConfigCh
		if !ok {
			return
		}

		log.Printf("[INFO]: Fetching url [%s]", feedConfig.Url)

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

		//log.Printf("[INFO]: %s: %d Kb", feedConfig.Url, len(feedConfig.Html)/1000)

		outCh <- feedConfig
	}
}
