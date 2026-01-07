package internal

import (
	"errors"
	"github.com/FG-GIS/boot-go-gator/internal/config"
	"log"
)

type state struct {
	configPointer *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error, not enough arguments.")
	}
	if len(cmd.args) > 1 {
		return errors.New("Error, too many arguments.")
	}
	err := s.configPointer.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	log.Println("GATOR -- User correctly set.")
	return nil
}
