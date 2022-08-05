package users

import (
	"context"
	"errors"

	"github.com/alesr/stdservices/users/internal/repository"
)

var _ repository.Repository = (*repositoryMock)(nil)

type repositoryMock struct {
	existsFunc        func(ctx context.Context, username, email string) (bool, error)
	insertFunc        func(ctx context.Context, user *repository.User) (*repository.User, error)
	selectByIDFunc    func(ctx context.Context, id string) (*repository.User, error)
	selectByEmailFunc func(ctx context.Context, email string) (*repository.User, error)
}

func (m *repositoryMock) Exists(ctx context.Context, username, email string) (bool, error) {
	if m.existsFunc == nil {
		return false, errors.New("repositoryMock.existsFunc is nil")
	}
	return m.existsFunc(ctx, username, email)
}

func (m *repositoryMock) Insert(ctx context.Context, user *repository.User) (*repository.User, error) {
	if m.insertFunc == nil {
		return nil, errors.New("repositoryMock.insertFunc is nil")
	}
	return m.insertFunc(ctx, user)
}

func (m *repositoryMock) SelectByID(ctx context.Context, id string) (*repository.User, error) {
	if m.selectByIDFunc == nil {
		return nil, errors.New("repositoryMock.selectByIDFunc is nil")
	}
	return m.selectByIDFunc(ctx, id)
}

func (m *repositoryMock) SelectByEmail(ctx context.Context, email string) (*repository.User, error) {
	if m.selectByEmailFunc == nil {
		return nil, errors.New("repositoryMock.selectByEmailFunc is nil")
	}
	return m.selectByEmailFunc(ctx, email)
}
