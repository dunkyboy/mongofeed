package main

import (
	"net/http"
	html_template "html/template"
	text_template "text/template"
	"fmt"
	"gopkg.in/mgo.v2"
	"flag"
	"os"
	"strings"
	"gopkg.in/mgo.v2/bson"
)

var mongoUrl string
var mongo *mgo.Session

func main() {
	fmt.Println("Starting mongofeed...")

	parseFlags()
	connectMongo()

	http.HandleFunc("/html/", htmlHandler)
	http.HandleFunc("/json/", jsonHandler)
	http.HandleFunc("/rss/", rssHandler)
	http.ListenAndServe(":18080", nil)
}

func parseFlags() {
	mongoUrl = *flag.String("m", "localhost:27017", "Mongo connection string")

	flag.Parse()
}

func connectMongo() {
	var mongoErr error;
	mongo, mongoErr = mgo.Dial(mongoUrl)
	if mongoErr != nil {
		fmt.Errorf("Error connecting to Mongo with connection string %s; exiting", mongoUrl)
		os.Exit(1)
	}
	fmt.Println("Connected to Mongo at", mongoUrl)
}

type Feed struct {
	Version string
	Title string
	HomepageUrl string
	FeedUrl string
	Items string
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := html_template.ParseFiles("templates/feed.html")
	t.Execute(w, "Antarctica")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path", r.URL.Path)

	path := strings.Split(r.URL.Path, "/")

	// path = "/json/db/collection"
	// path[0] = ""
	// path[1] = "json"
	db         := path[2]
	collection := path[3]

	feed, err := getFeed(r, db, collection)
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

func getFeed(r *http.Request, db string, collection string) (*Feed, error) {

	fmt.Printf("Fetching DB %v, collection %v\n", db, collection)

	var result []bson.M
	mongo.DB(db).C(collection).Find(bson.M{}).Sort("-$natural").Limit(1000).All(&result)

	var items string
	itemsBytes, err := bson.MarshalJSON(result)
	if err != nil {
		items = err.Error()
	} else {
		items = string(itemsBytes)
	}

	//items := make([]string, len(result))
	//for i, doc := range result {
	//	docStr, err := bson.MarshalJSON(doc)
	//	if err != nil {
	//		items[i] = err.Error()
	//	} else {
	//		items[i] = string(docStr[:])
	//	}
	//}

	fmt.Println("items:", items)

	return &Feed{
		"https://jsonfeed.org/version/1",
		fmt.Sprintf("mongofeed - %s", mongoUrl),
		r.Host,
		r.Host + r.RequestURI,
		items,
	}, nil
}
