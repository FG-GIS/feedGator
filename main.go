package main

import (
	"errors"
	"log"
	"os"

	"github.com/FG-GIS/boot-go-gator/internal/config"
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

func (c *commands) run(s *state, cmd command) error {
	err := c.commandList[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandList[name] = f
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %v", err)
	}
	s := state{
		configPointer: &conf,
	}
	c := commands{
		commandList: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalf("Error, not enough arguments.")
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	err = c.run(&s, cmd)
	if err != nil {
		log.Fatalf("GATOR Error: %v", err)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error, not enough arguments, Username is required for login.")
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
