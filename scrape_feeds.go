package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jradziejewski/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		ID: feedToFetch.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), params)

	for _, item := range feed.Channel.Item {
		fmt.Printf("- %s\n", item.Title)
	}

	return nil
}
