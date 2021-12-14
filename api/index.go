package handler

import (
	"io/ioutil"
	"net/http"
	"sgg-feed/feed"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// VercelのServerlessはGoを使うとローカルのファイルを読み込めないっぽい?
	response, err := http.Get("https://raw.githubusercontent.com/sugokunaritai-gakusei-group/sgg-feed/main/data.json")
	feed.ErrorHandling(err)
	defer response.Body.Close()
	rawJSON, err := ioutil.ReadAll(response.Body)
	feed.ErrorHandling(err)
	combinedFeedItems := feed.GetFeeds(rawJSON)
	readerArray, createdAt := feed.GenerateFeed(combinedFeedItems)
	feed.HostFeeds(readerArray, rawJSON, createdAt, w, r)
}
