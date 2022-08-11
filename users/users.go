package users

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/alesr/stdservices/pkg/validate"
	"github.com/alesr/stdservices/users/repository"
	"go.uber.org/zap"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	subject string = "%s: account verification"
	body    string = "Please click the following link to verify your account: %s/%s"
)

var _ Service = (*DefaultService)(nil)

type (

	// Service defines the service interface
	Service interface {
		// Create creates a new user and returns the created user with its ID and "user" role
		Create(ctx context.Context, in CreateUserInput) (*User, error)

		// Delete soft deletes a user by id
		Delete(ctx context.Context, id string) error

		// FetchByID fetches a non-deleted user by id and returns the user
		FetchByID(ctx context.Context, id string) (*User, error)

		// GenerateToken generates a JWT token for the user
		GenerateToken(ctx context.Context, email, password string) (string, error)

		// VerifyToken verifies a JWT token and returns the user username, id and role
		VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error)

		// SendEmailVerification sends an email verification to the user.
		// The user must be created before calling this method.
		SendEmailVerification(ctx context.Context, userID, to string) error
	}

	repo interface {
		Insert(ctx context.Context, user *repository.User) (*repository.User, error)
		SelectByID(ctx context.Context, id string) (*repository.User, error)
		SelectByEmail(ctx context.Context, email string) (*repository.User, error)
		DeleteByID(ctx context.Context, id string) error
		InsertEmailVerification(ctx context.Context, in repository.EmailVerification) error
	}

	emailer interface {
		Send(to, subject, body string) error
	}

	DefaultService struct {
		logger                    *zap.Logger
		appName                   string
		jwtSigningKey             string
		emailVerificationSecret   string
		emailVerificationEndpoint url.URL
		emailer                   emailer
		repo                      repo
	}
)

// New instantiates a new users service
func New(
	logger *zap.Logger,
	appName string,
	jwtSigningKey string,
	emailVerificationSecret string,
	emailVerificationEndpoint url.URL,
	emailer emailer,
	repo repo,
) *DefaultService {
	return &DefaultService{
		logger:                    logger,
		appName:                   appName,
		jwtSigningKey:             jwtSigningKey,
		emailVerificationSecret:   emailVerificationSecret,
		emailVerificationEndpoint: emailVerificationEndpoint,
		emailer:                   emailer,
		repo:                      repo,
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

	userID := uuid.NewString()

	insertedUser, err := s.repo.Insert(ctx, &repository.User{
		ID:            userID,
		Fullname:      in.Fullname,
		Username:      in.Username,
		Birthdate:     in.Birthdate,
		Email:         in.Email,
		EmailVerified: false,
		PasswordHash:  string(hash),
		Role:          string(in.Role),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateRecord) {
			return nil, errAlreadyExists
		}
		return nil, fmt.Errorf("could not insert user: %s", err)
	}

	if err := s.SendEmailVerification(ctx, userID, in.Email); err != nil {
		// It doesn't matter if the email verification fails, the user is created
		// and the user can request a new email verification
		s.logger.Error("could not send email verification", zap.String("user_id", userID), zap.Error(err))
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

func (s *DefaultService) SendEmailVerification(ctx context.Context, userID, to string) error {
	subject := fmt.Sprintf(subject, s.appName)

	// Generate verification token from user email and secret
	token, err := bcrypt.GenerateFromPassword(
		[]byte(to+s.emailVerificationSecret),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("could not email verification token: %s", err)
	}

	in := repository.EmailVerification{
		Token:     string(token),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24),
	}

	if err := s.repo.InsertEmailVerification(ctx, in); err != nil {
		return fmt.Errorf("could not insert email verification: %s", err)
	}

	body := fmt.Sprintf(body, s.emailVerificationEndpoint.String(), token)

	if err := s.emailer.Send(to, subject, body); err != nil {
		return fmt.Errorf("could not send email: %s", err)
	}
	return nil
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
		ID:            user.ID,
		Fullname:      user.Fullname,
		Username:      user.Username,
		Birthdate:     user.Birthdate,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Role:          role,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}
