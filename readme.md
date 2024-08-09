# Repository Pattern - Small Example

To demonstrate, in a small way, how to implement the repository pattern in Go, we’ll follow these steps:

1. Define the `UserRepository` interface.
2. Create a mock implementation of `UserRepository`.
3. Create a `UserService` that accepts an instance of the interface as a field
4. Write a unit test using the mock repository.

## 1. Define the UserRepository Interface

First, define the UserRepository *interface* that describes the methods for interacting with user data:

```go
package repository

type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository interface {
    FindUserByID(id int) (*User, error)
    SaveUser(user *User) error
}
```

Our interface is a *contract*. We are saying that any concrete instance that implements the interface, conforms to the contract.

This means we can *swap any implementation* of the interface, without changing any of the code. 

As an example, if you had a instance that connected to a PostgresDB, the implementation inside the methods `FindUserById` and `SaveUser` would have some sort of connection to the DB, and then some SQL to perform the action.

But with the repository pattern, we can create many instances that use different databases (or none at all) so that we are not tightly coupled to our database.

> :exclamation: Think of it this way - with the repository pattern we don't have to have our entire application running on our machine to test it. We can use in-memory DB, or test just our API layer, or our service layer, without having to load absolutely everything up, seed data into our DB, start the UI etc.

## 2. Create a Mock Implementation of UserRepository

Now, create a mock implementation of the UserRepository interface. This mock will be used in unit tests to *simulate* database interactions:

```go
package repository

import "errors"

type MockUserRepository struct {
    Users map[int]*User
    Err   error
}

func (m *MockUserRepository) FindUserByID(id int) (*User, error) {
    if m.Err != nil {
        return nil, m.Err
    }
    user, exists := m.Users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

func (m *MockUserRepository) SaveUser(user *User) error {
    if m.Err != nil {
        return m.Err
    }
    m.Users[user.ID] = user
    return nil
}
```

## 3. Write a Service That Uses The Interface

Of course, we now have the data access layer defined, but now we need to have some way of using that. We'll create a new folder - `service` - and add our `UserService`:

```go
package service

import "myapp/repository"

type UserService struct {
    Repo repository.UserRepository
}

func (s *UserService) GetUser(id int) (*repository.User, error) {
    return s.Repo.FindUserByID(id)
}

func (s *UserService) CreateUser(user *repository.User) error {
    return s.Repo.SaveUser(user)
}
```

The important bit here is our struct, `UserService`. This has a field `Repo` which can be set to *any implementation of the UserRepository interface*, and then uses that implementation to get and set data (using the `FindByUserId` and `SaveUser` methods).

Okay, maybe that doesn't make sense yet. Let's write a test that uses our UserService and the MockUserRepository to explain a bit more.

## 4. Write a Unit Test Using the Mock Repository

Next, we use the mock repository in a unit test. We're going to create an instance of our mock repository, and then when we define our service we'll set the Repo field to be our mock implementation.

```go
package service

import (
    "testing"
    "myapp/repository"
    "github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
    
    // Setup mock repository
    mockRepo := &repository.MockUserRepository{
        Users: map[int]*repository.User{
            1: {ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
        },
    }
    
    // Define the UserService with Repo set to our mockRepo instance 
    service := &UserService{Repo: mockRepo}
    
    //Now when we call our service.Getuser method, it will utilise the MockUserRepository to retrieve the user data
    user, err := service.GetUser(1)
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    
    //We can also check for errors when getting a user that doesn't exist
    user, err = service.GetUser(2)
    assert.Error(t, err)
    assert.Nil(t, user)
}

func TestCreateUser(t *testing.T) {
    //Setup our mockRepo to be empty
    mockRepo := &repository.MockUserRepository{
        Users: map[int]*repository.User{},
    }
    
    // Define the UserService with Repo set to our mockRepo instance 
    service := &UserService{Repo: mockRepo}
    
    // Add a new user to the "database"
    user := &repository.User{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"}
    err := service.CreateUser(user)
    assert.NoError(t, err)
    
    // Check that we don't have any other users, because we shouldn't!
    savedUser, err := mockRepo.FindUserByID(2)
    assert.NoError(t, err)
    assert.NotNil(t, savedUser)
    assert.Equal(t, "Jane Doe", savedUser.Name)
}
```

---

## In Brief

`MockUserRepository` implements the `UserRepository` *interface*. It has a map to store users and a way to simulate errors.

Unit Tests use the mock repository to test the service layer's `GetUser` and `CreateUser` methods. This ensures the service logic works correctly without needing a real database.

Tests focus on the service logic without depending on *actual database interactions*.

We can simulate different scenarios (e.g., errors, missing data) by adjusting the mock's behaviour.

---

## So What?

Well, take a look at what a `main.go` file would look like:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "myapp/repository" 

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
```

Our *application* when run, will use the postgres instance. It needs the database to be **up and available**.

If we *hadn't* used the repository pattern, our testing would have required the database to be running, for us to add and remove users from the database, reset data each time - all of which is a massive overhead and makes our test slower and more complex.

So whilst our app uses Postgres when *running*, we can separate out our testing of the service discretely using our own "in-memory database".

---

## Running The Tests

1. From the terminal in the project root run `go mod tidy` to install all dependencies
2. Run all tests using `go test -v ./...`

> You can't run the actual project unless you have a Postgres instance running and change all the connection details