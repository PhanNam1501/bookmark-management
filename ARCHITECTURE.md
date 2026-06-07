# 🏗️ Architecture - Bookmark Management Project

## Tổng Quan

Dự án sử dụng **Clean Architecture** với các layer:

```
HTTP Request
    ↓
[Handler Layer] - Parse input, return response
    ↓
[Service Layer] - Business logic
    ↓
[Repository Layer] - Database access
    ↓
Database (PostgreSQL)
```

---

## 📁 Folder Structure

```
bookmark-management/
├── cmd/
│   ├── api/
│   │   └── main.go              # Entry point
│   └── infrastructure/
│       ├── bootstrap.go         # DI Container
│       ├── database.go          # DB initialization
│       ├── redis.go             # Redis setup
│       └── config.go            # Config loading
│
├── internal/
│   ├── model/                   # Data entities
│   │   ├── base.go              # ID, timestamps
│   │   ├── user.go              # User entity
│   │   └── bookmark.go          # Bookmark entity
│   │
│   ├── repository/              # Database access
│   │   ├── repo.go              # Interface definition
│   │   ├── user.go              # User CRUD
│   │   ├── create.go            # Create bookmark
│   │   └── query.go             # Query bookmarks
│   │
│   ├── service/                 # Business logic
│   │   ├── user.go              # User service
│   │   ├── login.go             # Login logic
│   │   └── bookmark/
│   │       ├── service.go
│   │       ├── create.go
│   │       └── query.go
│   │
│   ├── handler/                 # HTTP handlers
│   │   ├── user.go              # User endpoints
│   │   ├── login.go             # Login endpoint
│   │   └── bookmark/
│   │       ├── handler.go
│   │       ├── create.go
│   │       └── query.go
│   │
│   ├── api/
│   │   ├── api.go               # Routes & DI
│   │   ├── config.go            # API config
│   │   └── middlewares/
│   │       └── jwt.go           # JWT middleware
│   │
│   └── test/                    # Integration tests
│
├── pkg/
│   ├── sqldb/                   # Database utilities
│   │   ├── dbclient.go
│   │   ├── config.go
│   │   └── migration.go
│   │
│   ├── jwtutils/                # JWT utilities
│   │   ├── generator.go         # Token generation
│   │   └── verifier.go          # Token validation
│   │
│   ├── utils/                   # Helper functions
│   └── dbutils/                 # DB error handling
│
├── migrations/                  # SQL migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_bookmarks_table.up.sql
│   └── 000002_create_bookmarks_table.down.sql
│
├── docker-compose.yaml          # Docker setup
├── Dockerfile                   # Container image
├── go.mod                       # Module definition
├── go.sum                       # Dependencies lock
└── README.md                    # Project documentation
```

---

## 🔄 Data Flow

### 1. Login Flow
```
POST /users/login
  ↓
Handler.Login (internal/handler/login.go)
  - Parse: {username, password}
  ↓
Service.Login (internal/service/login.go)
  - Get user by username
  - Verify password (bcrypt)
  - Create JWT claims
  - Sign token (private key)
  ↓
Repository.GetUserByUsername (internal/repository/user.go)
  - Query: SELECT * FROM users WHERE username = ?
  ↓
Response: {token: "eyJ..."}
```

### 2. Protected Request Flow
```
GET /v1/bookmarks?page=1&limit=10
  Authorization: Bearer eyJ...
  ↓
JWT Middleware (internal/api/middlewares/jwt.go)
  - Extract token from header
  - Validate signature (public key)
  - Store claims in context
  ↓
Handler.GetBookmarks (internal/handler/bookmark/query.go)
  - Extract userID from claims
  - Parse query params
  ↓
Service.GetBookmarks (internal/service/bookmark/query.go)
  - Calculate pagination
  ↓
Repository.QueryBookmarks (internal/repository/query.go)
  - Query: SELECT * FROM bookmarks WHERE user_id = ? ...
  - Count: SELECT COUNT(*) FROM bookmarks WHERE user_id = ?
  ↓
Response: {data: [...], pagination: {...}}
```

---

## 🔑 Key Components

### Handler Layer
- **Responsibility**: Parse HTTP input, call service, format response
- **Files**: `internal/handler/**/*.go`
- **Pattern**: Interface + implementation per entity

### Service Layer
- **Responsibility**: Business logic, validation, orchestration
- **Files**: `internal/service/**/*.go`
- **Pattern**: Inject repository + utilities

### Repository Layer
- **Responsibility**: Database operations (CRUD)
- **Files**: `internal/repository/**/*.go`
- **Pattern**: Interface per entity, implementation with GORM

### Model Layer
- **Responsibility**: Define data entities
- **Files**: `internal/model/*.go`
- **Pattern**: Struct with GORM tags

---

## 🔐 JWT Authentication

### Setup
```bash
# Generate keys
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
```

### Flow
1. **Generation** (pkg/jwtutils/generator.go)
   - Sign claims with private key (RS256)
   - Return JWT string

2. **Validation** (pkg/jwtutils/verifier.go)
   - Parse token
   - Verify signature with public key
   - Return claims if valid

3. **Middleware** (internal/api/middlewares/jwt.go)
   - Extract token from Authorization header
   - Validate token
   - Store claims in context

---

## 💾 Database

### Initialization
```
Bootstrap → InitDatabase → CreateGORMClient → RunMigrations
```

### Migrations
- **Location**: `migrations/` folder
- **Format**: SQL files with .up.sql and .down.sql
- **Tool**: golang-migrate/migrate
- **Trigger**: Auto-run in InitDatabase

### Schema
- **Users**: username, email, password_hash
- **Bookmarks**: user_id (FK), url, description, code
- **Soft Delete**: deleted_at column for both tables

---

## 🔄 Dependency Injection

### Pattern
```go
// In cmd/infrastructure/bootstrap.go or internal/api/api.go
repo := repository.NewRepository(db)
svc := service.NewService(repo, jwtGen)
handler := handler.NewHandler(svc)
```

### Advantages
- Easy to test (mock dependencies)
- Decoupled components
- Single responsibility

---

## 📊 Interfaces

### Repository
```go
type Repository interface {
    CreateBookmark(ctx context.Context, b *model.Bookmark) error
    QueryBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error)
    GetBookmarkCounts(ctx context.Context, userID string) (int64, error)
}
```

### Service
```go
type Service interface {
    CreateBookmark(ctx context.Context, userid, url, description string) (*model.Bookmark, error)
    GetBookmarks(ctx context.Context, userID string, limit, page int) (*GetBookmarksResult, error)
}
```

### Handler
```go
type Handler interface {
    CreateBookmark(c *gin.Context)
    GetBookmarks(c *gin.Context)
}
```

---

## 🚀 Running the Project

```bash
# Start services
docker-compose up -d

# Run server
go run cmd/api/main.go

# Run tests
go test ./...

# View API docs
http://localhost:8080/swagger/index.html
```

---

## 📝 Adding New Feature

### Example: Add "Delete Bookmark"

**1. Repository**
```go
type Repository interface {
    // ...
    DeleteBookmark(ctx context.Context, bookmarkID, userID string) error
}

func (r *bookmarkRepo) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
    return r.db.WithContext(ctx).
        Where("id = ? AND user_id = ?", bookmarkID, userID).
        Delete(&model.Bookmark{}).Error
}
```

**2. Service**
```go
type Service interface {
    // ...
    DeleteBookmark(ctx context.Context, userID, bookmarkID string) error
}

func (s *bookmarkService) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
    return s.r.DeleteBookmark(ctx, bookmarkID, userID)
}
```

**3. Handler**
```go
type Handler interface {
    // ...
    DeleteBookmark(c *gin.Context)
}

func (h *handler) DeleteBookmark(c *gin.Context) {
    // Implementation
}
```

**4. Routes**
```go
privateRoutes.DELETE("/v1/bookmarks/:id", handlers.BookmarkHandler.DeleteBookmark)
```

---

## ✅ Best Practices

1. ✅ Use interfaces for repositories and services
2. ✅ Inject dependencies, don't create instances
3. ✅ Use context for request cancellation
4. ✅ Wrap database errors with custom error types
5. ✅ Validate input at handler level
6. ✅ Hash passwords with bcrypt
7. ✅ Use migrations for schema changes
8. ✅ Add indexes for frequently queried columns
9. ✅ Add Swagger documentation for all endpoints
10. ✅ Write unit tests for services

---

## 🔗 Related Files

- [JWT Setup Guide](JWT.md)
- [Database Setup Guide](DATABASE.md)
- [API Endpoints](API.md)
- [Testing Guide](TESTING.md)
