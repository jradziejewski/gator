package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeleteFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("The database has been reset")
	return nil
}
