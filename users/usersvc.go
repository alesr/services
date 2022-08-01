package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alesr/services/users/internal/repository"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Enumerate roles

	RoleAdmin role = iota + 1
	RoleUser
)

type (
	role uint8

	AuthData struct {
		ID, Username, Role string
	}

	Service struct {
		jwtSigningKey string
		repo          repository.Repository
	}
)

func New(jwtSigningKey string, repo repository.Repository) *Service {
	return &Service{
		jwtSigningKey: jwtSigningKey,
		repo:          repo,
	}
}

// Create creates a new user and returns the created user
func (s *Service) Create(ctx context.Context, input CreateUserInput) (*User, error) {
	if err := input.validate(); err != nil {
		return nil, fmt.Errorf("could not validate create user input: %w", err)
	}

	exists, err := s.repo.Exists(ctx, input.Username, input.Email)
	if err != nil {
		return nil, fmt.Errorf("could not check if user exists: %s", err)
	}

	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}

	insertedUser, err := s.repo.Insert(ctx, &repository.User{
		ID:        uuid.NewString(),
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Username:  input.Username,
		Birthdate: input.Birthdate,
		Email:     input.Email,
		Hash:      string(hash),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("could not insert user: %s", err)
	}

	createdUser, err := newUserFromRepository(insertedUser)
	if err != nil {
		return nil, fmt.Errorf("could not create user from repository: %s", err)
	}

	return createdUser, nil
}

// FetchByID fetches a user by id and returns the user
func (s *Service) FetchByID(ctx context.Context, id string) (*User, error) {
	storageUser, err := s.repo.SelectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not select user by id: %s", err)
	}

	if storageUser == nil {
		return nil, ErrUserNotFound
	}

	user, err := newUserFromRepository(storageUser)
	if err != nil {
		return nil, fmt.Errorf("could not create user from repository: %s", err)
	}
	return user, nil
}

// AuthUser authenticates a user and returns the user token
func (s *Service) Auth(ctx context.Context, input AuthUserInput) (string, error) {
	if err := input.validate(); err != nil {
		return "", fmt.Errorf("could not validate auth user input: %w", err)
	}

	storageUser, err := s.repo.SelectByEmail(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("could not select user by email: %s", err)
	}

	if storageUser == nil {
		return "", ErrUserNotFound
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(storageUser.Hash), []byte(input.Password)); err != nil {
		return "", ErrPasswordInvalid
	}

	// Generate JWT
	token, err := s.generateJWT(storageUser.ID, RoleUser)
	if err != nil {
		return "", fmt.Errorf("could not generate JWT: %s", err)
	}
	return token, nil
}

func (s *Service) AccessToken(ctx context.Context, email, password string) (string, error) {
	// Fetch user by username
	storageUser, err := s.repo.SelectByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("could not select user by email: %s", err)
	}

	// Check if user exists
	if storageUser == nil {
		return "", ErrUserNotFound
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(storageUser.Hash), []byte(password)); err != nil {
		return "", ErrPasswordInvalid
	}

	// Generate JWT
	token, err := s.generateJWT(storageUser.ID, RoleUser)
	if err != nil {
		return "", fmt.Errorf("could not generate jwt: %s", err)
	}
	return token, nil
}

func (s *Service) Authorize(ctx context.Context, token string) (*AuthData, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSigningKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %s", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not parse token claims")
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("could not parse token expiration")
	}

	if exp < float64(time.Now().Unix()) {
		return nil, ErrTokenExpired
	}

	id, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("could not parse token id")
	}

	storageUser, err := s.repo.SelectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not select user by id: %s", err)
	}

	if storageUser == nil {
		return nil, ErrUserNotFound
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("could not parse token role")
	}

	return &AuthData{
		ID:       storageUser.ID,
		Username: storageUser.Username,
		Role:     role,
	}, nil
}

func (s *Service) generateJWT(userID string, role role) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": string(role),
		"exp":  time.Now().Add(time.Hour * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(s.jwtSigningKey))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %s", err)
	}
	return signedString, nil
}

func newUserFromRepository(user *repository.User) (*User, error) {
	var role role
	switch user.Role {
	case "user":
		role = RoleUser
	case "admin":
		role = RoleAdmin
	default:
		return nil, fmt.Errorf("invalid role: %s", user.Role)
	}

	return &User{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Birthdate: user.Birthdate,
		Email:     user.Email,
		Role:      role,
		CreatedAt: user.CreatedAt,
	}, nil
}
