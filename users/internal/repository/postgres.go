package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	// Enumerate postgresql query strings

	insertQuery        string = "INSERT INTO users (id,firstname,lastname,username,birthdate,email,hash,role,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
	selectByIDQuery    string = "SELECT id,firstname,lastname,username,birthdate,email,hash,role,created_at FROM users WHERE id = $1;"
	selectByEmailQuery string = "SELECT id,firstname,lastname,username,birthdate,email,hash,role,created_at FROM users WHERE email = $1;"
	existsQuery        string = "SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2;"
)

type (
	// Postgres represents a user repository instance with the given database connection
	Postgres struct{ *sqlx.DB }
)

// New creates a new user repository instance
func NewPostgres(dbConn *sqlx.DB) *Postgres {
	return &Postgres{dbConn}
}

// Insert inserts a new user into the database and returns the inserted user
func (p *Postgres) Insert(ctx context.Context, u *User) (*User, error) {
	insertStmt, err := p.PrepareContext(ctx, insertQuery)
	if err != nil {
		return nil, fmt.Errorf("could not prepare insert statement: %w", err)
	}
	defer insertStmt.Close()

	result, err := insertStmt.ExecContext(
		ctx, u.ID, u.Firstname, u.Lastname, u.Username,
		u.Birthdate, u.Email, u.Hash, u.Role, u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not insert user: %s", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %s", err)
	}

	if rowsAffected != 1 {
		return nil, fmt.Errorf("expected 1 row affected, got %d", rowsAffected)
	}
	return u, nil
}

func (p *Postgres) Exists(ctx context.Context, username, email string) (bool, error) {
	var count int
	if err := p.QueryRowContext(ctx, existsQuery, username, email).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return count > 0, nil
}

// SelectByID selects a user by id and returns the user
func (p *Postgres) SelectByID(ctx context.Context, id string) (*User, error) {
	user, err := p.selectUser(ctx, selectByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("could not select user by id: %s", err)
	}
	return user, nil
}

func (p *Postgres) SelectByEmail(ctx context.Context, email string) (*User, error) {
	user, err := p.selectUser(ctx, selectByEmailQuery, email)
	if err != nil {
		return nil, fmt.Errorf("could not select user by email: %s", err)
	}
	return user, nil
}

// selectUser executes the given query and returns the user
func (p *Postgres) selectUser(ctx context.Context, query, arg string) (*User, error) {
	var u User
	if err := p.QueryRowContext(ctx, query, arg).Scan(
		&u.ID, &u.Firstname, &u.Lastname, &u.Username,
		&u.Birthdate, &u.Email, &u.Hash, u.Role, &u.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select user: %s", err)
	}
	return &u, nil
}
