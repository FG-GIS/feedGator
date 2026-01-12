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
	s := getState()
	c := commands{
		commandList: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerGetUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerFeed))
	c.register("feeds", handlerGetFeeds)
	c.register("follow", middlewareLoggedIn(handlerSetFollow))
	c.register("following", middlewareLoggedIn(handlerShowFollowingUser))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	db, err := sql.Open("postgres", s.cfg.DBURL)

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
