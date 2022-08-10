package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alesr/stdservices/pkg/validate"
	"github.com/alesr/stdservices/users/internal/repository"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (

	// Service defines the service interface
	Service interface {
		Create(ctx context.Context, in CreateUserInput) (*User, error)
		Delete(ctx context.Context, id string) error
		FetchByID(ctx context.Context, id string) (*User, error)
		GenerateToken(ctx context.Context, email, password string) (string, error)
		VerifyToken(ctx context.Context, token string) (*User, error)
	}

	repo interface {
		Insert(ctx context.Context, user *repository.User) (*repository.User, error)
		SelectByID(ctx context.Context, id string) (*repository.User, error)
		SelectByEmail(ctx context.Context, email string) (*repository.User, error)
		DeleteByID(ctx context.Context, id string) error
	}

	DefaultService struct {
		jwtSigningKey string
		repo          repo
	}
)

// New instantiates a new users service
func New(jwtSigningKey string, repo repo) *DefaultService {
	return &DefaultService{
		jwtSigningKey: jwtSigningKey,
		repo:          repo,
	}
}

// Create creates a new user and returns the created user
func (s *DefaultService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	if err := in.validate(); err != nil {
		return nil, fmt.Errorf("could not validate create user input: %w", err)
	}

	// Block creation of users with admin role
	if in.Role == RoleAdmin {
		return nil, errCannotCreateAdminUser
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}

	insertedUser, err := s.repo.Insert(ctx, &repository.User{
		ID:           uuid.NewString(),
		Fullname:     in.Fullname,
		Username:     in.Username,
		Birthdate:    in.Birthdate,
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         string(in.Role),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateRecord) {
			return nil, errAlreadyExists
		}
		return nil, fmt.Errorf("could not insert user: %s", err)
	}

	createdUser, err := newUserFromRepository(insertedUser)
	if err != nil {
		return nil, fmt.Errorf("could not parse storage user to domain model: %s", err)
	}

	return createdUser, nil
}

// FetchByID fetches a user by id and returns the user
func (s *DefaultService) FetchByID(ctx context.Context, id string) (*User, error) {
	if err := validate.ID(id); err != nil {
		return nil, fmt.Errorf("could not validate id: %w", err)
	}

	storageUser, err := s.repo.SelectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not select user by id: %s", err)
	}

	if storageUser == nil {
		return nil, errNotFound
	}

	user, err := newUserFromRepository(storageUser)
	if err != nil {
		return nil, fmt.Errorf("could not parse storage user to domain model: %s", err)
	}
	return user, nil
}

func (s *DefaultService) Delete(ctx context.Context, id string) error {
	if err := validate.ID(id); err != nil {
		return fmt.Errorf("could not validate id: %w", err)
	}

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete user by id: %s", err)
	}
	return nil
}

// GenerateToken generates a JWT token for the user
func (s *DefaultService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	if err := validate.Email(email); err != nil {
		return "", fmt.Errorf("could not validate email: %s", err)
	}

	if err := validate.Password(password); err != nil {
		return "", fmt.Errorf("could not validate password: %s", err)
	}

	// Fetch user by username
	storageUser, err := s.repo.SelectByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("could not select user by email: %s", err)
	}

	// Check if user exists
	if storageUser == nil {
		return "", errNotFound
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(storageUser.PasswordHash), []byte(password)); err != nil {
		return "", errPasswordInvalid
	}

	// Generate JWT
	token, err := s.generateJWT(storageUser.ID, role(storageUser.Role))
	if err != nil {
		return "", fmt.Errorf("could not generate jwt: %s", err)
	}
	return token, nil
}

// VerifyToken verifies a JWT token and returns the authentication data
func (s *DefaultService) VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error) {
	if token == "" {
		return nil, errTokenEmpty
	}

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
		return nil, errTokenInvalid
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
		return nil, errTokenExpired
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
		return nil, errNotFound
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("could not parse token role")
	}

	return &VerifyTokenResponse{
		ID:       storageUser.ID,
		Username: storageUser.Username,
		Role:     role,
	}, nil
}

func (s *DefaultService) generateJWT(userID string, role role) (string, error) {
	if err := validate.ID(userID); err != nil {
		return "", fmt.Errorf("could not validate id: %w", err)
	}

	if err := role.validate(); err != nil {
		return "", errRoleInvalid
	}

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
		Fullname:  user.Fullname,
		Username:  user.Username,
		Birthdate: user.Birthdate,
		Email:     user.Email,
		Role:      role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
