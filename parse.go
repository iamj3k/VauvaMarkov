package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	url     string
	ok      bool
	topic   string
	thread  string
	comment string
)

// GetTopics Get list of topics
func GetTopics() []string {
	var topics []string
	resp, err := http.Get("https://www.vauva.fi/keskustelu/alue/aihe_vapaa")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".field-content").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if ok {
			topic = fmt.Sprintf("https://www.vauva.fi%s", href)
			topics = append(topics, topic)
		}
	})

	return topics
}

// GetThreads Gets discussion threads
func GetThreads(topic string) []string {
	var threads []string
	resp, err := http.Get(topic)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".views-field-title").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a").Attr("href")
		if ok {
			thread = fmt.Sprintf("https://www.vauva.fi%s", href)
			threads = append(threads, thread)
		}
	})
	return threads
}

// GetComments Make HTTP Request
func GetComments(url string) []string {
	var comments []string
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".field-item").Each(func(i int, s *goquery.Selection) {
		comment = s.Find("p").Text()
		comments = append(comments, comment)
	})
	return comments
}
