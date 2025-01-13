package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Accepts exactly two arguments <name> <url>")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("User not found in database")
	}

	now := time.Now().UTC()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	_, err = s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %s created successfully\n", cmd.args[0])
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		var authorName string
		if feed.AuthorName.Valid {
			authorName = feed.AuthorName.String
		}

		fmt.Printf("- '%s (%s) [author: %s]", feed.Name, feed.Url, authorName)
		fmt.Println()
	}

	return nil
}
