package main

import (
	"fmt"
	"os"

	"github.com/jradziejewski/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := &state{cfg: &cfg}

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	args = args[1:]

	cmds := newCommands()
	cmds.register("login", handlerLogin)

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
	cfg *config.Config
}

type command struct {
	name string
	args []string
}
type commands struct {
	handlers map[string]func(*state, command) error
}

// Functions

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Accepts exactly one argument <username>")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("Username has been set")
	return nil
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
