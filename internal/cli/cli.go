package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/rss"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("please enter a username")
	}

	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("please enter an unregistered username")
	}

	name := cmd.Args[0]

	user, err := s.DB.CreateUser(
		context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
		},
	)

	if err != nil {
		return err
	}

	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("User logged in: %s\n", name)
	fmt.Printf("User created: %+v\n", user)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if err := s.DB.ResetUser(context.Background()); err != nil {
		return err
	}

	fmt.Println("Reset database.")

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())

	if err != nil {
		return err
	}

	curUsr := s.Config.Username

	for _, user := range users {
		if user.Name == curUsr {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.DB.CreateFeed(
		context.Background(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Feed created:\n")
	fmt.Printf("ID: %v\nName: %s\nURL: %s\nUserID: %v\n", feed.ID, feed.Name, feed.Url, feed.UserID)

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.DB.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return err
	}

	fmt.Printf("Feedname: %s\nUsername: %s\n", follow.Feedname, follow.Username)

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\nURL: %s\nCreated by: %s\n", feed.Feedname, feed.Feedurl, feed.Username)
	}
	return nil
}

func HandlerFeedFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("usage: follow <url>")
	}

	url := cmd.Args[0]

	feed, err := s.DB.GetFeedUrl(context.Background(), url)
	if err != nil {
		return errors.New("feed does not seem to exist")
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.DB.CreateFeedFollow(context.Background(), params)

	if err != nil {
		return err
	}

	fmt.Printf("Feedname: %s\nUsername: %s\n", follow.Feedname, follow.Username)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.Feedname)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	url := cmd.Args[0]

	feed, err := s.DB.GetFeedUrl(context.Background(), url)
	if err != nil {
		return errors.New("feed does not seem to exist")
	}

	params := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.DB.UnfollowFeed(context.Background(), params); err != nil {
		return err
	}

	fmt.Printf("Unfollowed feed: %s",feed.Name)
	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.Username)
		if err != nil {
			return errors.New("current user does not exist")
		}
		return handler(s, cmd, user)
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("Command does not exist yet: %s", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(s *State, cmd Command) error) {
	c.Commands[name] = f
}
