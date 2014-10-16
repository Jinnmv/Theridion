package main

import (
  "net/http"
  "errors"
)

func Loader (out chan string, feedConfig FeedConfig) error {

  resp, err := http.Get(feedConfig.Url)
  if err!=nil {
    return error
  }
  defer resp.Body.Close()
  
  if resp.StatusCode < 200 && resp.StatusCode >= 400 {
    return errors.Error("Http request status", resp.StatusCode)
  }
  
  body, err := ioutil.ReadAll(resp.Body)
  if err!=nil {
    return err
  }
  
  out <- body
  
  return nil

}
