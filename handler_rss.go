package main

import (
	"context"
	"fmt"
)

var aggURL = "https://www.wagslane.dev/index.xml"

func aggFetch(s *state, cmd command) error {
	context := context.Background()
	feed, err := fetchFeed(context, aggURL)
	if err != nil {
		return fmt.Errorf("error retreving rss feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	return RSSFeed{}, nil
}
