package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/FG-GIS/feedGator/internal/config"
	"github.com/FG-GIS/feedGator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func scrapeFeeds(s *state, verbose bool) error {
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
	for _, item := range rss.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println("GATOR -- Error parsing publish time.")
			fmt.Println(item.PubDate)
			return err
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		})
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code != "23505" {
				fmt.Printf("GATOR -- Error saving post:\n -- %v\n", err)
			} else {
				continue
			}
		}
		if verbose {
			fmt.Printf("GATOR -- Post added => %v\n", post.Title)
		}
	}
	return nil
}

func printPost(post database.Post) {
	fmt.Println("------------------------------------------")
	fmt.Printf("* %s\n", post.Title)
	fmt.Printf("* %s\n", post.PublishedAt.Format(time.DateTime))
	fmt.Printf("***\n %s\n***\n", post.Description.String)
	fmt.Println("------------------------------------------")
}
