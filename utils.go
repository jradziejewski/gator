package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jradziejewski/gator/internal/config"
	"github.com/jradziejewski/gator/internal/database"
)

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

func parsePubDate(date string) (time.Time, error) {
	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700", // Common RSS
		"2006-01-02T15:04:05Z",            // ISO 8601
		"02 Jan 2006 15:04:05 -0700",      // Alternate format
	}
	var pubDate time.Time
	var err error

	for _, layout := range layouts {
		pubDate, err = time.Parse(layout, date)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, err
	}

	return pubDate, nil
}

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}
