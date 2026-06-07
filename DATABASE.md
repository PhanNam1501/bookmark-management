# 💾 Database Setup Guide

## Overview

This project uses:
- **Database**: PostgreSQL
- **ORM**: GORM
- **Migrations**: golang-migrate
- **Pattern**: Repository pattern with clean architecture

---

## Quick Start

### 1. Start PostgreSQL (Docker)

```bash
docker-compose up -d postgres
```

### 2. Run Migrations

Migrations run automatically when the app starts via `cmd/infrastructure/bootstrap.go`.

### 3. Verify Connection

```bash
# Connect to PostgreSQL
psql -U admin -d bookmark_service -h localhost

# List tables
\dt

# Exit
\q
```

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

### Bookmarks Table

```sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    url VARCHAR(2048) NOT NULL,
    description VARCHAR(255),
    code VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_created_at ON bookmarks(created_at);
```

---

## Migrations

### Location
```
migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_bookmarks_table.up.sql
└── 000002_create_bookmarks_table.down.sql
```

### File Structure

**Up (apply migration):**
```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
    -- ...
);
```

**Down (rollback migration):**
```sql
-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### Running Migrations

**Auto (on app start):**
```go
// cmd/infrastructure/database.go
func InitDatabase() (*gorm.DB, error) {
    db, err := sqldb.NewClient("")
    if err != nil {
        return nil, err
    }

    // Runs migrations automatically
    if err := sqldb.MigartePostgresDB(db, migrationPath, "up", 0); err != nil {
        return nil, err
    }

    return db, nil
}
```

**Manual:**
```bash
# Forward (apply)
migrate -path migrations -database "postgres://admin:admin@localhost:5432/bookmark_service" up

# Backward (rollback all)
migrate -path migrations -database "postgres://admin:admin@localhost:5432/bookmark_service" down

# Rollback 1 step
migrate -path migrations -database "postgres://admin:admin@localhost:5432/bookmark_service" down 1

# Check current version
migrate -path migrations -database "postgres://admin:admin@localhost:5432/bookmark_service" version
```

---

## Models (GORM)

### Base Model
```go
// internal/model/base.go
type Base struct {
    ID        string     `gorm:"type:uuid;primaryKey;column:id" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"-"`  // Soft delete
}
```

### User Model
```go
// internal/model/user.go
type User struct {
    Base
    Username    string `gorm:"column:username;uniqueIndex" json:"username"`
    Password    string `gorm:"column:password" json:"-"`
    Email       string `gorm:"column:email;uniqueIndex" json:"email"`
    DisplayName string `gorm:"column:display_name" json:"display_name"`
}
```

### Bookmark Model
```go
// internal/model/bookmark.go
type Bookmark struct {
    Base
    UserID      string `gorm:"column:user_id;index" json:"user_id"`
    URL         string `gorm:"column:url" json:"url"`
    Description string `gorm:"column:description" json:"description"`
    Code        string `gorm:"column:code;uniqueIndex" json:"code"`
}
```

---

## Repository Pattern

### Interface Definition

```go
// internal/repository/repo.go
type Repository interface {
    CreateBookmark(ctx context.Context, b *model.Bookmark) error
    QueryBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error)
    GetBookmarkCounts(ctx context.Context, userID string) (int64, error)
}
```

### Implementation

**Create:**
```go
// internal/repository/create.go
func (r *bookmarkRepo) CreateBookmark(ctx context.Context, b *model.Bookmark) error {
    return r.db.WithContext(ctx).Create(b).Error
}
```

**Query with Pagination:**
```go
// internal/repository/query.go
func (r *bookmarkRepo) QueryBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error) {
    var bookmarks []*model.Bookmark
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("created_at DESC").
        Limit(limit).
        Offset(offset).
        Find(&bookmarks).Error
    
    return bookmarks, err
}
```

**Count:**
```go
func (r *bookmarkRepo) GetBookmarkCounts(ctx context.Context, userID string) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&model.Bookmark{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    
    if err != nil {
        return 0, err
    }
    return count, nil
}
```

---

## Configuration

### Environment Variables

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=bookmark_service
DB_SSL_MODE=disable
```

### Load Configuration

```go
// pkg/sqldb/config.go
type config struct {
    Host     string `default:"localhost" envconfig:"DB_HOST"`
    User     string `default:"admin" envconfig:"DB_USER"`
    Password string `default:"admin" envconfig:"DB_PASSWORD"`
    DBName   string `default:"bookmark_service" envconfig:"DB_NAME"`
    Port     string `default:"5432" envconfig:"DB_PORT"`
    SSLMode  string `default:"disable" envconfig:"DB_SSL_MODE"`
}

func (cfg *config) GetDSN() string {
    return fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
    )
}
```

---

## GORM Operations

### Create
```go
bookmark := &model.Bookmark{
    UserID:      "user-id",
    URL:         "https://example.com",
    Description: "Example",
    Code:        "abc123",
}
db.Create(bookmark)
```

### Read
```go
// By ID
var bookmark model.Bookmark
db.First(&bookmark, "id = ?", id)

// Multiple with conditions
var bookmarks []model.Bookmark
db.Where("user_id = ?", userID).
    Order("created_at DESC").
    Limit(10).
    Offset(0).
    Find(&bookmarks)
```

### Update
```go
db.Model(&model.Bookmark{}).
    Where("id = ?", id).
    Update("description", "new description")
```

### Delete (Hard)
```go
db.Delete(&model.Bookmark{}, "id = ?", id)
```

### Delete (Soft)
```go
// Automatically sets deleted_at
db.Model(&model.Bookmark{}).Where("id = ?", id).Delete(&model.Bookmark{})

// Query excludes soft-deleted records by default
var bookmarks []model.Bookmark
db.Where("user_id = ?", userID).Find(&bookmarks)

// Include soft-deleted records
var bookmarks []model.Bookmark
db.Unscoped().Where("user_id = ?", userID).Find(&bookmarks)
```

---

## Useful Commands

### Connect to Database
```bash
psql -U admin -d bookmark_service -h localhost -p 5432
```

### View Table Structure
```sql
\d users
\d bookmarks
```

### View Indexes
```sql
SELECT * FROM pg_indexes WHERE tablename IN ('users', 'bookmarks');
```

### View All Records
```sql
SELECT * FROM users;
SELECT * FROM bookmarks;
```

### Count Records
```sql
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM bookmarks WHERE user_id = 'user-id';
```

### Check Foreign Keys
```sql
SELECT constraint_name, table_name, column_name
FROM information_schema.key_column_usage
WHERE table_name = 'bookmarks';
```

---

## Adding New Table

### Step 1: Create Migration Files

```bash
# Create migration files
touch migrations/000003_create_tags_table.up.sql
touch migrations/000003_create_tags_table.down.sql
```

### Step 2: Write SQL

```sql
-- migrations/000003_create_tags_table.up.sql
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

CREATE INDEX idx_tags_user_id ON tags(user_id);
```

```sql
-- migrations/000003_create_tags_table.down.sql
DROP TABLE IF EXISTS tags;
```

### Step 3: Restart App
```bash
go run cmd/api/main.go
```

Migration runs automatically!

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| "Connection refused" | Check if PostgreSQL is running: `docker-compose ps` |
| "Database does not exist" | Create database: `createdb -U admin bookmark_service` |
| "Permission denied" | Check user credentials in environment variables |
| "Migration failed" | Check SQL syntax, ensure table doesn't already exist |
| "Foreign key constraint failed" | Ensure referenced user exists before creating bookmark |

---

## Best Practices

1. ✅ **Migrations**: Use SQL files, version control them
2. ✅ **Indexes**: Add indexes for frequently queried columns
3. ✅ **Foreign Keys**: Define constraints in migrations
4. ✅ **Soft Delete**: Use deleted_at for archival
5. ✅ **Context**: Always use WithContext() for cancellation support
6. ✅ **Error Handling**: Wrap DB errors appropriately
7. ✅ **Pagination**: Always use LIMIT + OFFSET
8. ✅ **Query Optimization**: Use indexes, avoid N+1 queries

---

## Related Files

- [Architecture Overview](ARCHITECTURE.md)
- [JWT Guide](JWT.md)
- [API Endpoints](API.md)
