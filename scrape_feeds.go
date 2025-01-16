package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feed %s\n", feedToFetch.Name)

	for _, item := range feed.Channel.Item {
		pubDate, err := parsePubDate(item.PubDate)
		if err != nil {
			return err
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feedToFetch.ID,
		}
		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if isDuplicateError(err) {
				continue
			}
			return err
		}
		fmt.Printf("New post: '%s' added to database\n", post.Title)
	}

	return nil
}
