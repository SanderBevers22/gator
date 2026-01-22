package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gator/internal/cli"
	"gator/internal/config"
	"gator/internal/database"
)

const dbURL string = "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error while reading File: %v", err)
		log.Fatal(err)
	}

	state := &cli.State{
		DB:     dbQueries,
		Config: &cfg,
	}

	cmds := cli.Commands{
		Commands: make(map[string]func(*cli.State, cli.Command) error),
	}

	cmds.Register("login", cli.HandlerLogin)
	cmds.Register("register", cli.HandlerRegister)
	cmds.Register("reset", cli.HandlerReset)
	cmds.Register("users", cli.HandlerUsers)
	cmds.Register("agg", cli.HandlerAgg)
	cmds.Register("addfeed", cli.MiddlewareLoggedIn(cli.HandlerAddFeed))
	cmds.Register("feeds", cli.HandlerFeeds)
	cmds.Register("follow", cli.MiddlewareLoggedIn(cli.HandlerFeedFollow))
	cmds.Register("following", cli.MiddlewareLoggedIn(cli.HandlerFollowing))

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	cmd := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := cmds.Run(state, cmd); err != nil {
		log.Fatal(err)
	}
}
