package main

import (
	"net/http"
	html_template "html/template"
	text_template "text/template"
	"fmt"
	"gopkg.in/mgo.v2"
	"flag"
	"os"
)

var mongoUrl string
var mongo *mgo.Session

func main() {
	parseFlags()
	connectMongo()

	http.HandleFunc("/html", htmlHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/rss", rssHandler)
	http.ListenAndServe(":18080", nil)
}

func parseFlags() {
	mongoUrl = *flag.String("m", nil, "Mongo connection string")

	flag.Parse()
}

func connectMongo() {
	var mongoErr error;
	mongo, mongoErr = mgo.Dial(mongoUrl)
	if mongoErr != nil {
		fmt.Errorf("Error connecting to Mongo with connection string %s; exiting", mongoUrl)
		os.Exit(1)
	}
}

type Feed struct {
	Version string
	Title string
	HomepageUrl string
	FeedUrl string
	Items []string
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := html_template.ParseFiles("templates/feed.html")
	t.Execute(w, "Antarctica")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	feed, err := getFeed(r)
	if err != nil {
		fmt.Errorf("Error retrieving JSON feed: %v", err)
		return  // TODO
	}

	t, _ := text_template.ParseFiles("templates/feed.json")
	t.Execute(w, feed)
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func getFeed(r *http.Request) (*Feed, error) {

	r.

	mongo.DB()

	return &Feed{
		"https://jsonfeed.org/version/1",
		"MongoFeed",
		r.Host,
		r.Host + r.RequestURI,
		make([]string, 0),
	}, nil
}
