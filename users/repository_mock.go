package users

import (
	"context"
	"errors"

	"github.com/alesr/stdservices/internal/users/repository"
)

var _ repo = (*repositoryMock)(nil)

type repositoryMock struct {
	insertFunc        func(ctx context.Context, user *repository.User) (*repository.User, error)
	selectByIDFunc    func(ctx context.Context, id string) (*repository.User, error)
	selectByEmailFunc func(ctx context.Context, email string) (*repository.User, error)
	deleteByIDFunc    func(ctx context.Context, id string) error
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

func (m *repositoryMock) DeleteByID(ctx context.Context, id string) error {
	if m.deleteByIDFunc == nil {
		return errors.New("repositoryMock.deleteByIDfunc is nil")
	}
	return m.deleteByIDFunc(ctx, id)
}
