package repository

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
    DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) FindUserByID(id int) (*User, error) {
    var user User
    query := "SELECT id, name, email FROM users WHERE id = $1"
    row := r.DB.QueryRow(query, id)

    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    return &user, nil
}

func (r *PostgresUserRepository) SaveUser(user *User) error {
    query := `
    INSERT INTO users (name, email) 
    VALUES ($1, $2) 
    RETURNING id`
    
    err := r.DB.QueryRow(query, user.Name, user.Email).Scan(&user.ID)
    if err != nil {
        return err
    }
    
    return nil
}