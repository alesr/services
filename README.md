# stdservices
![master](https://github.com/alesr/stdservices/actions/workflows/ci.yaml/badge.svg)


This is a collection of small services that I use in my side projects.

This set of services includes a Docker Compose file with a Postgres database and migrations for integration testing.

If you decide to use these services, you should use the database migrations as a reference for your project.


```
------------------------------------------------------------------------
stdservices
------------------------------------------------------------------------
db-down                        remove the test database container and its volumes
db                             spins up the test database
lint                           run go format, vet and lint code
migrate                        executes the migrations towards the test database
psql                           executes a psql command to connect to the test database
test-it                        run integration tests
test-unit                      run unit tests
test                           Run unit and integration tests
```

## Testing

So far, this project has implemented unit, integration tests, and code linting.

To run unit tests, run `make test-unit`.

Integrations tests require connecting to a Postgres database and migrations found on the "migrations" folder.

The command `make test-integration` will spin up a Docker container with a Postgres database, execute the migrations, and run tests following the naming convention `TestIntegration...`. 

The command `make test` runs the linter, which includes `go fmt`, `go vet` and `statickcheck`, unit tests, and integration tests.

Finally, you can use the service mocks to unit test your application.

---
## users

The user service implements a set of CRUD operations for users. It is used in conjunction with JWT authentication and includes
a repository layer for storing users in a database. A PostgreSQL implementation is provided for convenience.

```go
type Service interface {
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
	}
```

### Upcoming features
    - Edit user
    - Password reset
    - Email verification
    - Feed service
    - Profile service
    ...
