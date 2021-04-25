package main

import (
	"io/ioutil"
	"sgg-feed/feed"
)

func main() {
	rawJSON, err := ioutil.ReadFile("./data.json")
	feed.ErrorHandling(err)
	combinedFeedItems := feed.GetFeeds(rawJSON)
	readerArray := feed.GenerateFeed(combinedFeedItems)
	feed.HostFeeds(readerArray)
}
