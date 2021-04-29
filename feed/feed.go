package feed

import (
	"encoding/json"
	"fmt"
	"github.com/p1ass/feeder"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

func crawlFeed(item feedType, feedChannel chan feeder.Crawler) {
	if item.Service == "qiita" {
		feedChannel <- feeder.NewAtomCrawler("https://qiita.com/" + item.Value + "/feed")
	}
	if item.Service == "zenn" {
		feedChannel <- feeder.NewRSSCrawler("https://zenn.dev/" + item.Value + "/feed")
	}
	if item.Service == "other" {
		feedChannel <- feeder.NewRSSCrawler(item.Value)
	}
}

func GetFeeds(rawJSON []byte) []*feeder.Item {
	var data []feedType
	var combinedItem []feeder.Crawler
	json.Unmarshal(rawJSON, &data)
	feedChannel := make(chan feeder.Crawler)
	defer close(feedChannel)
	for _, item := range data {
		go crawlFeed(item, feedChannel)
	}
	var count int
	for {
		select {
		case feed := <- feedChannel:
			combinedItem = append(combinedItem, feed)
			count++
			if count == len(data) {
				feedItem, err := feeder.Crawl(combinedItem...)
				ErrorHandling(err)
				return feedItem
			}
		}
	}
}

func GenerateFeed(combinedFeedItems []*feeder.Item) []*string {
	finalFeed := &feeder.Feed{
		Title: "SGG feed",
		Link: &feeder.Link{Href: "https:/sgg-feed.appspot.com/rss"},
		Description: "Integrated RSS&JSON feed of SGG Community. Articles are all written in Japanese.",
		Author: &feeder.Author{
			Name: "SGG Members",
		},
		Created: time.Now(),
		Items: combinedFeedItems,
	}
	rss, err := finalFeed.ToRSS()
	ErrorHandling(err)
	json, err := finalFeed.ToJSON()
	ErrorHandling(err)
	return []*string{ &rss, &json }
}

func HostFeeds(readerArray []*string) {
	rssReader := readerArray[0]
	jsonReader := readerArray[1]
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/rss", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/rss+xml;charset=UTF-8")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
		fmt.Println("Content Header has been set")
		if (*req).Method == "OPTIONS" {
			fmt.Println("Preflight")
			return
		}
		reader := strings.NewReader(*rssReader)
		_, err := io.Copy(writer, reader)
		ErrorHandling(err)
		return
	})
	http.HandleFunc("/api", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
		fmt.Println("Content Header has been set")
		if (*req).Method == "OPTIONS" {
			fmt.Println("Preflight")
			return
		}
		reader := strings.NewReader(*jsonReader)
		_, err := io.Copy(writer, reader)
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

