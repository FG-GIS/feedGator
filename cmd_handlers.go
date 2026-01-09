package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/FG-GIS/boot-go-gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error, not enough arguments, Username is required for login.")
	}
	if len(cmd.args) > 1 {
		return errors.New("Error, too many arguments.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		log.Fatalf("Error, non-existant username.\n %v", err)
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	log.Println("GATOR -- User correctly set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error, not enough arguments, Username is required for registration.")
	}
	if len(cmd.args) > 1 {
		return errors.New("Error, too many arguments.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		return errors.New("User already registered, exiting.")
	}
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        int32(uuid.New()[0]),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	log.Println("GATOR -- User correctly registered.")
	log.Printf("User generated: %v", usr)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("Error, too many arguments.")
	}
	err := s.db.Reset(context.Background())
	if err != nil {
		log.Println("Error re-setting table")
		return err
	}
	log.Println("GATOR -- Table cleansed")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("Error, too many arguments.")
	}
	usrSlice, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Println("Error gathering users from table")
		return err
	}
	for _, usr := range usrSlice {
		if usr == s.cfg.CurrentUserName {
			log.Printf("* %v (current)\n", usr)
		} else {
			log.Printf("* %v\n", usr)
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

}
