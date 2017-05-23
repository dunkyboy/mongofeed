package main

import (
	"net/http"
	html_template "html/template"
	text_template"text/template"
	"fmt"
)

func main() {
	http.HandleFunc("/html", htmlHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/rss", rssHandler)
	http.ListenAndServe(":18080", nil)
}

type Feed struct {
	Version string
	Title string
	HomepageUrl string
	FeedUrl string
	Items []string
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := html_template.ParseFiles("html/feed.html")
	t.Execute(w, "Antarctica")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	feed, err := getFeed(r)
	if err != nil {
		fmt.Errorf("Error retrieving JSON feed: %v", err)
		return  // TODO
	}

	t, _ := text_template.ParseFiles("html/feed.json")
	t.Execute(w, feed)
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func getFeed(r *http.Request) (f *Feed, err error) {
	return &Feed{
		"https://jsonfeed.org/version/1",
		"MongoFeed",
		r.Host,
		r.Host + r.RequestURI,
		make([]string, 0),
	}, nil
}
