package feed

import (
	"encoding/json"
	"fmt"
	"github.com/p1ass/feeder"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func ErrorHandling(errMessage error) {
	if errMessage != nil {
		log.Fatal(errMessage)
		os.Exit(1)
	}
}

type feedType struct {
	Service string `json:"service"`
	Value string `json:"value"`
}

type individualFeeds struct {
	Name string `json:"name"`
	Feeds []feedType `json:"feeds"`
}

func GetFeeds(rawJSON []byte) []*feeder.Item {
	var data []individualFeeds
	var combinedItem []feeder.Crawler
	json.Unmarshal(rawJSON, &data)
	for _, t := range data {
		for _, item := range t.Feeds {
			fmt.Println(t.Name, item.Service)
			var rss feeder.Crawler
			if item.Service == "qiita" {
				rss = feeder.NewAtomCrawler("https://qiita.com/" + item.Value + "/feed")
			}
			if item.Service == "zenn" {
				rss = feeder.NewRSSCrawler("https://zenn.dev/" + item.Value + "/feed")
			}
			if item.Service == "other" {
				rss = feeder.NewRSSCrawler(item.Value)
			}
			combinedItem = append(combinedItem, rss)
		}
	}
	feedItem, err := feeder.Crawl(combinedItem...)
	ErrorHandling(err)
	return feedItem
}

func GenerateFeed(combinedFeedItems []*feeder.Item) []io.Reader {
	finalFeed := &feeder.Feed{
		Title: "SGG feed",
		Link: &feeder.Link{Href: "https://example.com/feed"},
		Description: "Integrated RSS&JSON feed of SGG Community. Articles are all written in Japanese.",
		Author: &feeder.Author{
			Name: "SGG Members",
		},
		Created: time.Now(),
		Items: combinedFeedItems,
	}
	rss, err := finalFeed.ToRSSReader()
	ErrorHandling(err)
	json, err := finalFeed.ToJSONReader()
	ErrorHandling(err)
	return []io.Reader{ rss, json }
}

func HostFeeds(readerArray []io.Reader) {
	rssReader := &readerArray[0]
	jsonReader := &readerArray[1]
	http.HandleFunc("/rss", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/rss+xml")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
		if (*req).Method == "OPTIONS" {
			return
		}
		_, err := io.Copy(writer, *rssReader)
		ErrorHandling(err)
		return
	})
	http.HandleFunc("/api", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
		if (*req).Method == "OPTIONS" {
			return
		}
		_, err := io.Copy(writer, *jsonReader)
		ErrorHandling(err)
		return
	})
	portName := os.Getenv("PORT")
	if portName == "" {
		portName = "3432"
	}
	fmt.Println("RSS feed has been published at http://localhost:" + portName)
	err := http.ListenAndServe(":" + portName, nil)
	ErrorHandling(err)
}

