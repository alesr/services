package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	// Enumerate postgresql query strings

	insertQuery string = `INSERT INTO users (id,fullname,username,birthdate,email,password_hash,
	role,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING 
	id,fullname,username,birthdate,email,password_hash,role,created_at,updated_at;`

	selectByIDQuery string = `SELECT id,fullname,username,birthdate,email,password_hash,
	role,created_at,updated_at FROM users WHERE id = $1;`

	selectByEmailQuery string = `SELECT id,fullname,username,birthdate,email,password_hash,
	role,created_at,updated_at FROM users WHERE email = $1;`

	existsQuery string = "SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2;"
)

// Postgres represents a user repository instance with the given database connection
type Postgres struct{ *sqlx.DB }

// New creates a new user repository instance
func NewPostgres(dbConn *sqlx.DB) *Postgres {
	return &Postgres{dbConn}
}

func (p *Postgres) Insert(ctx context.Context, u *User) (*User, error) {
	var res User

	if err := p.QueryRowContext(
		ctx,
		insertQuery,
		u.ID,
		u.Fullname,
		u.Username,
		u.Birthdate,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(
		&res.ID, &res.Fullname, &res.Username, &res.Birthdate, &res.Email,
		&res.PasswordHash, &res.Role, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not scan inserted user: %s", err)
	}
	return &res, nil
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
		&u.ID, &u.Fullname, &u.Username, &u.Birthdate,
		&u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select user: %s", err)
	}
	return &u, nil
}
