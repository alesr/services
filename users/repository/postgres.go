package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

const (
	// Enumerate postgresql query strings

	insertQuery string = `INSERT INTO users (id,fullname,username,birthdate,email,email_verified,password_hash,
	role,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING 
	id,fullname,username,birthdate,email,email_verified,password_hash,role,created_at,updated_at;`

	selectByIDQuery string = `SELECT id,fullname,username,birthdate,email,email_verified,
	password_hash,role,created_at,updated_at FROM users WHERE id = $1 AND deleted_at IS NULL;`

	selectByEmailQuery string = `SELECT id,fullname,username,birthdate,email,email_verified,
	password_hash,role,created_at,updated_at FROM users WHERE email = $1 AND deleted_at IS NULL;`

	deleteByIDQuery string = "UPDATE users SET deleted_at = NOW() WHERE id = $1;"

	insertEmailVerificationQuery string = `INSERT INTO email_verifications 
	(code,user_id,created_at,expires_at) VALUES ($1,$2,$3,$4);`
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
		ctx, insertQuery, u.ID, u.Fullname, u.Username,
		u.Birthdate, u.Email, u.EmailVerified, u.PasswordHash,
		u.Role, u.CreatedAt, u.UpdatedAt,
	).Scan(
		&res.ID, &res.Fullname, &res.Username, &res.Birthdate, &res.Email,
		&res.EmailVerified, &res.PasswordHash, &res.Role, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return nil, ErrDuplicateRecord
		}
		return nil, fmt.Errorf("could not scan inserted user: %s", err)
	}
	return &res, nil
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
		&u.ID, &u.Fullname, &u.Username, &u.Birthdate, &u.Email,
		&u.EmailVerified, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("could not select user: %s", err)
	}
	return &u, nil
}

func (p *Postgres) DeleteByID(ctx context.Context, id string) error {
	res, err := p.ExecContext(ctx, deleteByIDQuery, id)
	if err != nil {
		return fmt.Errorf("could not delete user: %s", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %s", err)
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (p *Postgres) InsertEmailVerification(ctx context.Context, in EmailVerification) error {
	_, err := p.ExecContext(ctx, insertEmailVerificationQuery, in.Code, in.UserID, in.CreatedAt, in.ExpiresAt)
	if err != nil {
		return fmt.Errorf("could not insert email verification: %s", err)
	}
	return nil
}
