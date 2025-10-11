package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/simonjwhitlock/bootdev-go-blogaggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time interval>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time interval: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
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
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("failed to find feed to fetch: %w", err)
		return
	}
	fetchedTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	markFetchedParams := database.MarkFeedFetchedParams{
		ID:            feedToFetch.ID,
		LastFetchedAt: fetchedTime,
	}
	err = s.db.MarkFeedFetched(context.Background(), markFetchedParams)
	if err != nil {
		fmt.Println("failed to mark feed as fetched: %w", err)
		return
	}
	fetchedFeed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		fmt.Println("failed to fetch feed: %w", err)
		return
	}

	fmt.Printf("fetched: %v\n***********************\n", fetchedFeed.Channel.Title)
	for _, item := range fetchedFeed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("time format wrong: %v\n", err)
			pubTime = time.Now()
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: pubTime,
			FeedID:      feedToFetch.ID,
		}
		createdPost, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if !strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				log.Printf("failed to create post: %v\n", err)
			}
		} else {
			fmt.Printf("post created: %v\n", createdPost.Title)
		}
	}

	fmt.Printf("***********************\n\n")
}

func handlerGetUsersPosts(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		limitint, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = int32(limitint)

	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <number or posts (optional)>", cmd.Name)
	}
	getPostsParams := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostForUser(context.Background(), getPostsParams)
	if err != nil {
		return err
	}
	for _, fetchedPost := range posts {
		fmt.Printf("Feed Name: %v\n", fetchedPost.FeedName)
		fmt.Printf("Post Title: %v\n", fetchedPost.Title)
		fmt.Printf("Post URL: %v\n", fetchedPost.Url)
		fmt.Printf("Post publised at: %v\n", fetchedPost.PublishedAt)
		fmt.Printf("Post Description:\n %v\n", fetchedPost.Description)
		fmt.Println("*****************************************")
	}
	return nil
}
