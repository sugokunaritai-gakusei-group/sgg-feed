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

	ioutil.WriteFile("./static/feed.xml", []byte(*readerArray[0]), 0666)
	ioutil.WriteFile("./static/feed.json", []byte(*readerArray[1]), 0666)
}
