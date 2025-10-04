package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonjwhitlock/bootdev-go-blogaggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	context := context.Background()
	_, err := s.db.GetUser(context, name)
	if err != nil {
		return fmt.Errorf("user not found in database: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerCreateUser(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	context := context.Background()
	userUuid := uuid.New()
	userName := cmd.Args[0]
	time := time.Now()
	userArg := database.CreateUserParams{
		ID:        userUuid,
		CreatedAt: time,
		UpdatedAt: time,
		Name:      userName,
	}
	newUser, err := s.db.CreateUser(context, userArg)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Println(newUser)

	fmt.Println("created user: ", userName)

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully!")
	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	context := context.Background()

	err := s.db.ResetUsers(context)
	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	context := context.Background()

	users, err := s.db.GetUsers(context)
	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println(user.Name, "(current)")
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}
