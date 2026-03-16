# Simple REST API

A modern Go REST API for managing students, teachers, courses, and enrollments. Built with pure Go's `net/http`, featuring comprehensive testing, middleware support, and PostgreSQL database integration.

## Table of Contents

- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Libraries](#libraries)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Database Migrations](#database-migrations)
- [Handlers](#handlers)
- [Middlewares](#middlewares)
- [Repositories](#repositories)
- [Mocks](#mocks)
- [Logging](#logging)
- [Testing](#testing)

## Getting Started

### Prerequisites

- Go 1.25.1 or later
- PostgreSQL 12 or later
- Git

### Installation

```bash
git clone https://github.com/DenysShpak0116/simple-rest-api
cd simple-rest-api
go mod download
```

## Project Structure

```
simple-rest-api/
├── cmd/                          # Command-line applications
│   ├── migrator/main.go          # Database migration runner
│   └── server/main.go            # Server entry point
├── config/
│   └── config.yaml               # Configuration file
├── internal/
│   ├── config/
│   │   └── config.go             # Configuration loading
│   ├── database/
│   │   └── conn.go               # Database connection setup
│   ├── handlers/                 # HTTP request handlers
│   │   ├── courses/              # Course-related handlers & tests
│   │   │   ├── create.go         # Create course handler
│   │   │   ├── delete.go         # Delete course handler
│   │   │   ├── get.go            # Get course handler
│   │   │   ├── list.go           # List courses handler
│   │   │   ├── update.go         # Update course handler
│   │   │   ├── *_test.go         # Unit tests
│   │   │   └── mocks/            # Mock interfaces for testing
│   │   ├── students/             # Student-related handlers & tests
│   │   ├── teachers/             # Teacher-related handlers & tests
│   │   ├── enrollments/          # Enrollment-related handlers & tests
│   │   ├── responses/            # Response formatters
│   │   ├── types/                # Type definitions
│   │   ├── server.go             # Server setup and routing
│   │   ├── routes.go             # Route definitions
│   │   └── repositories_interfaces.go  # Repository interfaces
│   ├── models/                   # Data models
│   │   ├── course.go
│   │   ├── student.go
│   │   ├── teacher.go
│   │   └── enrollment.go
│   └── repository/               # Data access layer
│       ├── course_repo.go        # Course repository
│       ├── student_repo.go       # Student repository
│       ├── teacher_repo.go       # Teacher repository
│       ├── enrollment_repo.go    # Enrollment repository
│       └── errors.go             # Repository error definitions
├── pkg/
│   ├── http/
│   │   ├── middleware/
│   │   │   ├── middleware.go     # Middleware chaining utility
│   │   │   └── logger/
│   │   │       └── logger.go     # Request logging middleware
│   │   └── render/
│   │       └── render.go         # JSON response rendering
│   └── slogpretty/               # Pretty structured logging
├── migrations/                   # Database migrations (Goose)
│   ├── 20260313204945_init_tables.sql
│   └── 20260316103133_add_on_delete.sql
├── docs/                         # Swagger API documentation
├── tests/                        # Integration tests
├── go.mod                        # Dependencies
├── go.sum                        # Dependency checksums
└── README.md                     # This file
```

## Libraries

### Core Dependencies

| Library                                  | Purpose                          |
| ---------------------------------------- | -------------------------------- |
| `github.com/lib/pq`                      | PostgreSQL driver for Go         |
| `github.com/swaggo/swag`                 | Swagger documentation generation |
| `github.com/swaggo/http-swagger`         | Swagger UI for HTTP              |
| `github.com/pressly/goose/v3`            | Database migration tool          |
| `github.com/go-playground/validator/v10` | Input validation                 |
| `github.com/ilyakaznacheev/cleanenv`     | Configuration loading            |
| `github.com/fatih/color`                 | Colored terminal output          |

### Testing Dependencies

| Library                       | Purpose                                  |
| ----------------------------- | ---------------------------------------- |
| `github.com/stretchr/testify` | Testing assertions and mocking framework |

## Configuration

### Config File Format

Create `config/config.yaml`:

```yaml
env: "local"
connectionString: "postgres://user:password@localhost:5432/school_db?sslmode=disable"
http:
  host: "localhost"
  port: 8081
  timeout: 10s
```

### Configuration Fields

- **env**: Environment mode (`local`, `dev`, `prod`)
- **connectionString**: PostgreSQL connection string (DSN format)
- **http.host**: Server host address
- **http.port**: Server port
- **http.timeout**: Request timeout duration

### Loading Configuration

The application loads configuration in this order:

1. Command-line flag: `-config /path/to/config.yaml`
2. Environment variable: `CONFIG_PATH=/path/to/config.yaml`
3. Panics if not found

```bash
go run cmd/server/main.go -config config/config.yaml
```

## Running the Application

### Prerequisites

1. **Create PostgreSQL Database**

```bash
createdb school_db
```

2. **Create Configuration File**

3. **Run Migrations**

```bash
go run cmd/migrator/main.go -config config/config.yaml
```

### Start the Server

```bash
go run cmd/server/main.go -config config/config.yaml
```

Server will start on `http://localhost:8081` (depends on your configs)

### Check API Documentation

Navigate to:

```
http://localhost:8081/swagger/index.html
```

## Database Migrations

The project uses **Goose** for database migrations.

### Migration Files Location

```
migrations/
├── 20260313204945_init_tables.sql    # Initial schema
└── 20260316103133_add_on_delete.sql  # Schema modifications
```

### Migration Commands

```bash
# Run all pending migrations
go run cmd/migrator/main.go -config config/config.yaml

# Or use goose CLI directly
goose -dir migrations postgres "postgres://user:pass@localhost:5432/db" up
goose -dir migrations postgres "postgres://user:pass@localhost:5432/db" down
```

### Initial Schema (20260313204945_init_tables.sql)

Creates tables for:

- **students**: Student information with email uniqueness constraint
- **teachers**: Teacher information with department
- **courses**: Course information with teacher foreign key
- **enrollments**: Student-course many-to-many relationships

Includes indexes on frequently queried columns.

### Second Migration (20260316103133_add_on_delete.sql)

Modifies course-teacher relationship:

- Changes teacher_id to nullable
- Updates foreign key constraint from ON DELETE CASCADE to ON DELETE SET NULL

## Handlers

Handlers are the HTTP request processing functions. They handle request validation, business logic, and response formatting.

### Handler Structure

Each resource has 5 handlers:

```
<resource>/
├── create.go        # POST /api/v1/<resource>
├── get.go           # GET /api/v1/<resource>/{id}
├── list.go          # GET /api/v1/<resource>
├── update.go        # PUT /api/v1/<resource>/{id}
└── delete.go        # DELETE /api/v1/<resource>/{id}
```

### Handler Implementation Pattern

```go
type CourseCreator interface {
	CreateCourse(course *models.Course) (int, error)
}

func Create(logger *slog.Logger, creator CourseCreator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Decode and validate request body
		req, problems, err := render.DecodeValid[createCourseRequest](r)
		if len(problems) > 0 {
			responses.ValidationError(w, r, problems)
			return
		}

		// 2. Apply business logic
		id, err := creator.CreateCourse(&course)
		if err != nil {
			responses.Error(w, r, http.StatusInternalServerError, "failed to create")
			return
		}

		// 3. Render response
		render.Encode(w, r, http.StatusCreated, types.ApiResponse[courseResponse]{
			Message: "created successfully",
			Data:    &response,
		})
	})
}
```

### Handler Dependencies

Handlers accept interfaces, not concrete types:

- Loose coupling with repository layer
- Easy mocking for tests
- Implementation can be swapped

### Supported Resources

- **Students**: Create, Read, Update, Delete, List
- **Teachers**: Create, Read, Update, Delete, List
- **Courses**: Create, Read, Update, Delete, List
- **Enrollments**: Create, Delete

## Middlewares

Middlewares are functions that wrap HTTP handlers to add cross-cutting concerns.

### Middleware Chain

```go
handler = middleware.Chain(handler, loggerMiddleware, authMiddleware)
```

Applied in **reverse order** (bottom-up).

### Built-in Middlewares

#### 1. Logger Middleware

**Location**: `pkg/http/middleware/logger/logger.go`

**Purpose**: Logs HTTP request details

**What it logs**:

- HTTP method (GET, POST, etc.)
- Request path
- Request ID (from X-Request-ID header)

**Code**:

```go
func New(log *slog.Logger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			log.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", reqID),
			)
			next.ServeHTTP(w, r)
		})
	}
}
```

### Middleware Utilities

**Location**: `pkg/http/middleware/middleware.go`

**Middleware Type Definition**:

```go
type Middleware func(http.Handler) http.Handler
```

**Chain Function**:

```go
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
```

Applies middlewares in reverse order, so first middleware in list is outermost.

### Creating Custom Middleware

```go
func CustomMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Before
		log.Info("before request")

		// Call next handler
		next.ServeHTTP(w, r)

		// After
		log.Info("after request")
	})
}
```

## Repositories

Repositories encapsulate data access logic. They provide an abstraction over the database layer.

### Repository Pattern

```go
type StudentRepository interface {
	CreateStudent(student *models.Student) (int, error)
	GetStudentByID(id int) (*models.Student, error)
	GetAllStudents() ([]*models.Student, error)
	UpdateStudent(student *models.Student) error
	DeleteStudent(id int) error
}
```

### Repository Implementation

Each resource has a repository:

```
repository/
├── student_repo.go      # StudentRepository implementation
├── teacher_repo.go      # TeacherRepository implementation
├── course_repo.go       # CourseRepository implementation
├── enrollment_repo.go   # EnrollmentRepository implementation
└── errors.go            # Shared error definitions
```

### Common Repository Operations

```go
type StudentRepository struct {
	db *sql.DB
}

// Create
func (r *StudentRepository) CreateStudent(student *models.Student) (int, error) {
	var id int
	err := r.db.QueryRow(query, student.FirstName, student.LastName, student.Email).Scan(&id)
	return id, err
}

// Read
func (r *StudentRepository) GetStudentByID(id int) (*models.Student, error) {
	student := &models.Student{}
	err := r.db.QueryRow(query, id).Scan(&student.ID, &student.FirstName, ...)
	return student, err
}

// Update
func (r *StudentRepository) UpdateStudent(student *models.Student) error {
	result, err := r.db.Exec(query, student.FirstName, ...)
	return r.checkRowsAffected(result, err)
}

// Delete
func (r *StudentRepository) DeleteStudent(id int) error {
	result, err := r.db.Exec(query, id)
	return r.checkRowsAffected(result, err)
}

// List
func (r *StudentRepository) GetAllStudents() ([]*models.Student, error) {
	rows, err := r.db.Query(query)
	defer rows.Close()
	// Scan and build slice
}
```

### Repository Error Handling

**Location**: `internal/repository/errors.go`

```go
var (
	ErrStudentNotFound      = errors.New("student not found")
	ErrTeacherNotFound      = errors.New("teacher not found")
	ErrCourseNotFound       = errors.New("course not found")
	ErrEnrollmentNotFound   = errors.New("enrollment not found")
	ErrDuplicateEmail       = errors.New("email already exists")
	ErrEnrollmentConflict   = errors.New("student already enrolled in this course")
	ErrForeignKeyConstraint = errors.New("foreign key constraint violation")
)
```

**Error Handling in Handlers**:

```go
course, err := getter.GetCourseByID(courseID)
if err != nil {
	if errors.Is(err, repository.ErrCourseNotFound) {
		responses.Error(w, r, http.StatusNotFound, "course not found")
		return
	}
	responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve")
	return
}
```

## Mocks

Mocks are test doubles that simulate repository behavior without hitting the database.

### Mock Generation

Mocks are generated using **mockery** mock package:

```bash
mockery
```

### Mock Files Location

```
<resource>/mocks/
├── mock_StudentCreator.go
├── mock_StudentDeleter.go
├── mock_StudentGetter.go
├── mock_StudentsLister.go
└── mock_StudentUpdater.go
```

## Logging

The application uses **structured logging** with Go's `slog` package for consistent, machine-readable logs.

### Logging Configuration

**Location**: `cmd/server/main.go`

```go
func setupLogger(env string, w io.Writer) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(slogpretty.NewPrettyHandler(w))
	case envDev:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{...}))
	case envProd:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{...}))
	}
}
```

### Log Levels

- **DEBUG**: Detailed information for debugging (enabled in local/dev only)
- **INFO**: General information about application flow
- **WARN**: Warning messages for potentially problematic situations
- **ERROR**: Error messages for failures that need investigation

### Pretty Logging (Local Environment)

**Location**: `pkg/slogpretty/slogpretty.go`

Features:

- Color-coded output
- Human-readable format
- Perfect for development

**Output Example**:

```
INFO: request method=GET path=/api/v1/courses request_id=abc123def456
WARN: validation failed problems=[{Field: "email", Message: "invalid email format"}]
ERROR: database connection failed error="connection refused"
```

### Structured Logging Usage

```go
log.Info("course created successfully",
	slog.Int("id", courseID),
	slog.String("title", course.Title),
	slog.Int("teacher_id", *course.TeacherID),
)

log.Warn("student not found",
	slog.Int("student_id", studentID),
)

log.Error("database query failed",
	slog.String("error", err.Error()),
	slog.String("query", "SELECT * FROM students"),
)
```

## Testing

The project uses **table-driven tests** with **testify mocks** for comprehensive test coverage.

### Test Structure

```
<resource>/
├── create_test.go
├── delete_test.go
├── get_test.go
├── list_test.go
├── update_test.go
└── mocks/
```

### Test Pattern

```go
func TestCreateCourse(t *testing.T) {
	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		requestBody    map[string]any
		mockSetup      func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Course Created",
			requestBody: map[string]any{
				"title":       "Intro to Go",
				"description": "A comprehensive guide",
				"teacher_id":  1,
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
				validator.EXPECT().GetTeacherByID(1).Return(&models.Teacher{ID: 1}, nil)
				creator.EXPECT().CreateCourse(mock.Anything).Return(42, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":42`,
		},
		{
			name: "Error - Validation Failed",
			requestBody: map[string]any{
				"title": "A",
			},
			mockSetup: func(creator *mocks.MockCourseCreator, validator *mocks.MockTeacherValidator) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "validation failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCreator := mocks.NewMockCourseCreator(t)
			mockValidator := mocks.NewMockTeacherValidator(t)

			tc.mockSetup(mockCreator, mockValidator)

			handler := courses.Create(silentLogger, mockCreator, mockValidator)

			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedBody)
		})
	}
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/handlers/courses/...

# Run specific test
go test -run TestCreateCourse ./internal/handlers/courses

# Run with verbose output
go test -v ./internal/handlers/courses/...

# Run with coverage
go test -cover ./...
```

### Test Coverage

The project includes comprehensive tests for:

- **Create handlers**: Success cases, validation errors, business logic errors, database errors
- **Get handlers**: Success, invalid ID, not found, database error
- **List handlers**: Success with data, empty results, database error
- **Update handlers**: Success, validation, not found, database error
- **Delete handlers**: Success, invalid ID, not found, database error

### Test Utilities

- **httptest.NewRequest()**: Create test HTTP requests
- **httptest.NewRecorder()**: Capture response data
- **testify/assert**: Assertions for test conditions
- **testify/require**: Assertions that fail the test immediately
- **testify/mock**: Create and manage mock expectations
