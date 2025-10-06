package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonjwhitlock/bootdev-go-blogaggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	newFeed, err := s.db.AddFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Println("New feed added:", newFeed.Name)
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	fmt.Println("lists feeds")
	feedsList, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feedsList {
		fmt.Printf("Feed Name: %v\n", feed.Name)
		fmt.Printf("Feed URL: %v\n", feed.Url)
		fmt.Printf("Feed User: %v\n\n", feed.UserName)
	}

	return nil
}
