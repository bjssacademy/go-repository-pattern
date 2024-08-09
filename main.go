package main

import (
	"database/sql"
	"fmt"
	"gorepository/repository" // Adjust the import path as needed
	"log"

	_ "github.com/lib/pq"
)

func main() {
    connStr := "user=youruser dbname=yourdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    userRepo := repository.NewPostgresUserRepository(db)

    // Create a new user
    newUser := &repository.User{Name: "Alice", Email: "alice@example.com"}
    err = userRepo.SaveUser(newUser)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("New user ID: %d\n", newUser.ID)

    // Retrieve a user by ID
    user, err := userRepo.FindUserByID(newUser.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User found: %s, %s\n", user.Name, user.Email)
}