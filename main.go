package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/FG-GIS/boot-go-gator/internal/config"
	"github.com/FG-GIS/boot-go-gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %v", err)
	}
	s := state{
		cfg: &conf,
	}
	c := commands{
		commandList: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)

	db, err := sql.Open("postgres", conf.DBURL)

	dbQueries := database.New(db)
	s.db = dbQueries

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
