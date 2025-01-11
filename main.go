package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jradziejewski/gator/internal/config"
	"github.com/jradziejewski/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := &state{cfg: &cfg}

	db, err := sql.Open("postgres", st.cfg.DBUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	st.db = dbQueries

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	args = args[1:]

	cmds := newCommands()
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	cmd := command{
		name: args[0],
		args: args[1:],
	}

	if err := cmds.run(st, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Struct Types

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}
type commands struct {
	handlers map[string]func(*state, command) error
}

// Handlers

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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("The users table has been reset")
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

// Util Functions

func newCommands() *commands {
	return &commands{
		handlers: make(map[string]func(*state, command) error),
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}
