package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/FG-GIS/boot-go-gator/internal/config"
	"github.com/FG-GIS/boot-go-gator/internal/database"
)

func getState() state {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %v", err)
	}
	s := state{
		cfg: &conf,
	}
	return s

}

func printFeeds(feedSlice []database.GetFeedsRow) {
	fmt.Printf("GATOR -- Printing out feeds (%d)", len(feedSlice))
	for _, f := range feedSlice {
		fmt.Printf("* ID:							%s\n", f.ID)
		fmt.Printf("* Name:							%s\n", f.Name)
		fmt.Printf("* URL:							%s\n", f.Url)
		fmt.Printf("* :							%s\n", f.User)
	}
}

func printFollowing(followSlice []database.GetFeedFollowsForUserRow, user string) {
	fmt.Printf("GATOR -- User ==> %s - is following:\n", user)
	if len(followSlice) == 0 {
		fmt.Printf("(empty)")
	}
	for _, follow := range followSlice {
		fmt.Printf("* %s\n", follow.FeedName)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUsr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			fmt.Println("GATOR -- Error getting current user.")
			return err
		}
		err = handler(s, cmd, currentUsr)
		if err != nil {
			return err
		}
		return nil
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("GATOR -- Error getting next feed to fetch.")
		return err
	}
	updateParams := database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = s.db.MarkFeedFetched(context.Background(), updateParams)
	if err != nil {
		fmt.Println("GATOR -- Error updating fetch time.")
		return err
	}
	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Println("GATOR -- Error fetching feed.")
		return err
	}
	rss.printTitles()
	return nil
}
