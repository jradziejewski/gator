package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Accepts exactly one argument <url>")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Feed with given link not found in database")
	}

	now := time.Now().UTC()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %s followed feed '%s'\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range following {
		fmt.Printf("- '%s'\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Accepts exactly one argument <url>")
	}

	params := database.DeleteFollowParams{
		UserID: user.ID,
		Url:    cmd.args[0],
	}

	err := s.db.DeleteFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Unfollowed successfully")

	return nil
}
