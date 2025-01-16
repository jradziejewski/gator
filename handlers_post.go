package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jradziejewski/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("Accepts max one argument <limit>")
	}

	var limit int

	if len(cmd.args) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	} else {
		limit = 2
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("- '%s (%s) ", post.Title, post.Url)
		fmt.Println()
	}

	return nil
}
