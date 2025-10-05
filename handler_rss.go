package main

import (
	"context"
	"fmt"
)

var aggURL = "https://www.wagslane.dev/index.xml"

func handlerGetRssFeed(s *state, cmd command) error {
	context := context.Background()
	feed, err := rss.fetchFeed(context, aggURL)
	if err != nil {
		return fmt.Errorf("error retreving rss feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}
