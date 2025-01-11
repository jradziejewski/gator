package main

import (
	"database/sql"
	"fmt"
	"os"

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
