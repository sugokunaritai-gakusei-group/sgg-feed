package feed

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/p1ass/feeder"
)

func ErrorHandling(errMessage error) {
	if errMessage != nil {
		log.Println(errMessage)
	}
}

type feedType struct {
	Service string `json:"service"`
	Value   string `json:"value"`
}

func GetFeeds(rawJSON []byte) []*feeder.Item {
	var data []feedType
	var feedItems []*feeder.Item
	json.Unmarshal(rawJSON, &data)
	for _, item := range data {
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
		feedItem, err := feeder.Crawl(rss)
		ErrorHandling(err)
		feedItems = append(feedItems, feedItem...)
	}
	return feedItems
}

func GenerateFeed(combinedFeedItems []*feeder.Item) ([]*string, time.Time) {
	finalFeed := &feeder.Feed{
		Title:       "SGG feed",
		Link:        &feeder.Link{Href: "sgg-feed.appspot.com/"},
		Description: "Integrated RSS&JSON feed of SGG Community. Articles are all written in Japanese.",
		Author: &feeder.Author{
			Name: "SGG Members",
		},
		Created: time.Now(),
		Items:   combinedFeedItems,
	}
	rss, err := finalFeed.ToRSS()
	ErrorHandling(err)
	json, err := finalFeed.ToJSON()
	ErrorHandling(err)
	return []*string{&rss, &json}, time.Now()
}

func HostFeeds(readerArray []*string, rawJSON []byte, createdAt time.Time) {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/rss", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/rss+xml;charset=UTF-8")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		fmt.Println("Content Header has been set")
		if (*req).Method == "OPTIONS" {
			fmt.Println("Preflight")
			return
		}

		elapsedTime := time.Since(createdAt)
		if elapsedTime.Hours()*60+elapsedTime.Minutes() > 60 {
			readerArray, createdAt = GenerateFeed(GetFeeds(rawJSON))
			println("regenerated")
		}
		rssReader := readerArray[0]

		reader := strings.NewReader(*rssReader)
		_, err := io.Copy(writer, reader)
		ErrorHandling(err)
		return
	})
	http.HandleFunc("/api", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		fmt.Println("Content Header has been set")
		if (*req).Method == "OPTIONS" {
			fmt.Println("Preflight")
			return
		}

		elapsedTime := time.Since(createdAt)
		if elapsedTime.Hours()*60+elapsedTime.Minutes() > 60 {
			readerArray, createdAt = GenerateFeed(GetFeeds(rawJSON))
		}
		jsonReader := readerArray[1]

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
	err := http.ListenAndServe(":"+portName, nil)
	ErrorHandling(err)
}
