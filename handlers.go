package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/gator/internal/database"
)

// User handlers

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Accepts exactly one argument <username>")
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User %s does not exist in database", username)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("Username has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Accepts exactly one argument <username>")
	}

	now := time.Now().UTC()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
	}
	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s created successfully\n", cmd.args[0])
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("- %s", user)
		if s.cfg.CurrentUserName == user {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}

	return nil
}

// Feed handlers

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Accepts exactly two arguments <name> <url>")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
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

// Reset

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("The users table has been reset")
	return nil
}

// Agg

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v", feed)

	return nil
}
