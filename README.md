# stdservices
![Coverage](https://img.shields.io/badge/Coverage-71.8%25-brightgreen)
![master](https://github.com/alesr/stdservices/actions/workflows/ci.yaml/badge.svg)


This is a collection of small services that I use in my side projects.

## users

The user service implements a set of CRUD operations for users. It is used in conjunction with JWT authentication and includes
a repository layer for storing users in a database. A PostgreSQL implementation is provided for convenience.

```go
Service interface {
    Create(ctx context.Context, in CreateUserInput) (*User, error)
    FetchByID(ctx context.Context, id string) (*User, error)
    GenerateToken(ctx context.Context, email, password string) (string, error)
    VerifyToken(ctx context.Context, token string) (*User, error)
}
```

### Upcoming features
    - Edit user
    - Soft delete
    - Password reset
    - Email verification
