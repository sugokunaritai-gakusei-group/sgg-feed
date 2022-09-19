package main

import (
	"io/ioutil"
	"sgg-feed/feed"
)

func main() {
	rawJSON, err := ioutil.ReadFile("./data.json")
	feed.ErrorHandling(err)
	combinedFeedItems := feed.GetFeeds(rawJSON)
	readerArray, _ := feed.GenerateFeed(combinedFeedItems)

	ioutil.WriteFile("./static/rss", []byte(*readerArray[0]), 0666)
	ioutil.WriteFile("./static/api", []byte(*readerArray[1]), 0666)
}
