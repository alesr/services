package users

import (
	"context"
	"errors"
)

var _ Service = (*MockService)(nil)

type MockService struct {
	CreateFunc                func(ctx context.Context, in CreateUserInput) (*User, error)
	DeleteFunc                func(ctx context.Context, id string) error
	FetchByIDFunc             func(ctx context.Context, id string) (*User, error)
	GenerateTokenFunc         func(ctx context.Context, email, password string) (string, error)
	VerifyTokenFunc           func(ctx context.Context, token string) (*VerifyTokenResponse, error)
	SendEmailVerificationFunc func(ctx context.Context, userID, to string) error
}

func (m *MockService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
	if m.CreateFunc == nil {
		return nil, errors.New("MockService.CreateFunc is nil")
	}
	return m.CreateFunc(ctx, in)
}

func (m *MockService) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return errors.New("MockService.DeleteFunc is nil")
	}
	return m.DeleteFunc(ctx, id)
}

func (m *MockService) FetchByID(ctx context.Context, id string) (*User, error) {
	if m.FetchByIDFunc == nil {
		return nil, errors.New("MockService.FetchByIDFunc is nil")
	}
	return m.FetchByIDFunc(ctx, id)
}

func (m *MockService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	if m.GenerateTokenFunc == nil {
		return "", errors.New("MockService.GenerateTokenFunc is nil")
	}
	return m.GenerateTokenFunc(ctx, email, password)
}

func (m *MockService) VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error) {
	if m.VerifyTokenFunc == nil {
		return nil, errors.New("MockService.VerifyTokenFunc is nil")
	}
	return m.VerifyTokenFunc(ctx, token)
}

func (m *MockService) SendEmailVerification(ctx context.Context, userID, to string) error {
	if m.SendEmailVerificationFunc == nil {
		return errors.New("MockService.SendEmailVerificationFunc is nil")
	}
	return m.SendEmailVerificationFunc(ctx, userID, to)
}
