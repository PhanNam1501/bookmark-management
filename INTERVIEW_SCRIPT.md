# 🎯 Bookmark Management API - Interview Script

## 📋 Project Overview

**Project Name:** Bookmark Management API  
**Duration:** Self-learning project  
**Role:** Full-stack Backend Developer  
**Status:** Production-ready  

Tôi xây dựng một **REST API backend** cho ứng dụng quản lý bookmarks với các tính năng xác thực người dùng, quản lý dữ liệu, và tối ưu hóa performance.

---

## 🏗️ Architecture Overview

### 1. **Architecture Pattern: Clean Architecture**

Dự án sử dụng **Clean Architecture** để đảm bảo:
- 📁 **Separation of Concerns** - mỗi layer có trách nhiệm riêng
- 🔄 **Dependency Inversion** - các layer không phụ thuộc vào chi tiết implement
- ✅ **Testability** - dễ viết unit tests bằng cách mock dependencies

**4 Main Layers:**
```
┌─────────────────────────────────────┐
│     HTTP Handler Layer              │ (Gin framework - parse request, return response)
├─────────────────────────────────────┤
│     Service Layer                   │ (Business logic - validation, orchestration)
├─────────────────────────────────────┤
│     Repository Layer                │ (Database operations - CRUD)
├─────────────────────────────────────┤
│     Database Layer (PostgreSQL)      │ (Persistent data storage)
└─────────────────────────────────────┘
```

**Why Clean Architecture?**
- Dễ dàng thêm tính năng mới mà không ảnh hưởng đến code cũ
- Nếu cần thay đổi database từ PostgreSQL sang MySQL, chỉ cần sửa Repository layer
- Unit tests dễ viết vì có thể mock Repository

---

## 🛠️ Tech Stack

### **Backend Framework & Language**
| Technology | Version | Purpose |
|------------|---------|---------|
| **Go (Golang)** | 1.25.0 | Backend language - fast, concurrent, simple |
| **Gin Gonic** | v1.12.0 | Web framework - HTTP routing, middleware |

**Tại sao Go?**
- ✅ Fast execution & low latency
- ✅ Excellent concurrency support (goroutines)
- ✅ Single binary output - dễ deploy
- ✅ Strong standard library
- ✅ Phù hợp cho microservices

**Tại sao Gin?**
- ✅ Lightweight & fast routing
- ✅ Built-in middleware system
- ✅ Easy to define routes and handlers
- ✅ Good error handling

---

### **Database & ORM**
| Technology | Version | Purpose |
|------------|---------|---------|
| **PostgreSQL** | Latest | Relational database - reliable, ACID-compliant |
| **GORM** | v1.31.1 | ORM - type-safe queries, migrations |
| **golang-migrate** | v4.19.1 | Database migrations - version control for schema |

**Database Schema:**
```sql
Users Table:
├── id (UUID)
├── username (UNIQUE)
├── email
├── password_hash (bcrypt)
├── display_name
├── created_at / updated_at / deleted_at (soft delete)

Bookmarks Table:
├── id (UUID)
├── user_id (FK → Users.id)
├── url
├── description
├── code (shortened URL identifier)
├── created_at / updated_at / deleted_at
```

**Why GORM?**
- ✅ Type-safe queries
- ✅ Automatic migration support
- ✅ Relationship handling
- ✅ Query builder pattern

**Why golang-migrate?**
- ✅ Version control for database schema
- ✅ Up/Down migrations for rollback
- ✅ Team collaboration (tracked in git)

---

### **Authentication & Security**
| Technology | Purpose |
|------------|---------|
| **JWT (JSON Web Tokens)** | Stateless authentication |
| **RS256 (RSA)** | Asymmetric cryptography - private/public key pairs |
| **Bcrypt** | Password hashing - secure, salted hashing |

**JWT Flow:**
```
1. User Login
   POST /users/login → username + password
   
2. Server Validation
   ✓ Query user from DB
   ✓ Compare password with bcrypt
   
3. Token Generation
   ✓ Create JWT claims: {sub: user_id, iat, exp}
   ✓ Sign with private.pem (RS256)
   ✓ Return token string
   
4. Protected Request
   GET /v1/bookmarks
   Header: Authorization: Bearer <token>
   
5. Middleware Validation
   ✓ Extract token from header
   ✓ Verify signature with public.pem
   ✓ Extract user_id from claims
   ✓ Query bookmarks for that user
```

**Security Measures:**
- ✅ Passwords hashed with bcrypt (one-way)
- ✅ JWT signed with RS256 (cannot be forged without private key)
- ✅ Token expiration (24 hours)
- ✅ Soft deletes (data not permanently removed)
- ✅ User isolation (each user only sees their own data)

---

### **Caching & Performance**
| Technology | Version | Purpose |
|------------|---------|---------|
| **Redis** | v9.18.0 | In-memory cache - fast data access |

**Use Cases:**
- ✅ Rate limiting (prevent API abuse)
- ✅ Session caching (improve login performance)
- ✅ Query result caching (reduce DB load)

---

### **Testing**
| Technology | Version | Purpose |
|------------|---------|---------|
| **Testify** | v1.11.1 | Unit testing framework - assertions, mocking |
| **Miniredis** | v2.37.0 | Mock Redis - test without real Redis server |

**Testing Strategy:**
- Unit tests for services (business logic)
- Integration tests for handlers (with mock database)
- Mock repositories to isolate tests
- 80%+ code coverage target

---

### **API Documentation**
| Technology | Purpose |
|------------|---------|
| **Swagger/OpenAPI** | Auto-generate interactive API docs |
| **swaggo/gin-swagger** | Integrate Swagger into Gin routes |

**Benefits:**
- ✅ Auto-generated from code comments
- ✅ Interactive UI at `/swagger/index.html`
- ✅ Can test endpoints directly in browser
- ✅ Serves as documentation

---

### **Logging & Monitoring**
| Technology | Version | Purpose |
|------------|---------|---------|
| **Zerolog** | v1.35.0 | Structured logging - JSON format, fast |

**Log Levels:**
```
DEBUG → INFO → WARN → ERROR → FATAL
```

**Benefits:**
- ✅ Structured JSON logs (easy to parse & analyze)
- ✅ Fast performance (low overhead)
- ✅ Easy integration with ELK, Datadog, etc.

---

### **Environment & Configuration**
| Technology | Version | Purpose |
|------------|---------|---------|
| **envconfig** | v1.4.0 | Load config from environment variables |

**Why envconfig?**
- ✅ 12-factor app compliance
- ✅ Easy to manage secrets (API keys, DB credentials)
- ✅ Different configs for dev/staging/production
- ✅ No need to commit sensitive data to git

**Example:**
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
REDIS_URL=redis://localhost:6379
JWT_PRIVATE_KEY_PATH=/path/to/private.pem
```

---

## 🚀 Core Features

### **1. User Management**
- ✅ User registration (username, email, password)
- ✅ Password hashing with bcrypt
- ✅ User login with JWT token generation
- ✅ Get current user info from token

**Endpoints:**
```
POST   /users/register       → Create new user
POST   /users/login          → Login & get JWT token
GET    /v1/self/info         → Get current user (protected)
```

---

### **2. Bookmark Management**
- ✅ Create bookmarks (URL + description)
- ✅ Query bookmarks with pagination
- ✅ Soft delete bookmarks
- ✅ User-specific data isolation

**Endpoints:**
```
POST   /v1/bookmarks         → Create bookmark (protected)
GET    /v1/bookmarks?page=1&limit=10  → List bookmarks (protected)
DELETE /v1/bookmarks/:id     → Delete bookmark (protected)
```

**Example Response:**
```json
{
  "id": "f1e27b78-97a5-4456-8163-6a83fade5dab",
  "user_id": "user-123",
  "url": "https://github.com",
  "description": "GitHub - Where world builds software",
  "code": "f94128e1032209442b3e1bf6e8ef60",
  "created_at": "2026-06-03T00:00:00Z"
}
```

---

### **3. Advanced Features (Additional)**
- ✅ URL shortening (convert long URLs to short codes)
- ✅ URL redirect service (short code → original URL)
- ✅ Rate limiting (prevent API abuse using Redis)
- ✅ CSV import (bulk import bookmarks)
- ✅ Password management service
- ✅ Link management & analytics

---

## 📊 Database Design

### **Users Table**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP (soft delete)
);
```

### **Bookmarks Table**
```sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    url VARCHAR(2048),
    description TEXT,
    code VARCHAR(32) UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP (soft delete)
);

CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
```

**Indexing Strategy:**
- ✅ Primary key on `id` (automatic)
- ✅ Unique index on `username` (fast login lookups)
- ✅ Foreign key index on `user_id` (fast bookmark queries)
- ✅ Unique index on `code` (fast URL redirect lookups)

---

## 🐳 Deployment & DevOps

### **Docker & Containerization**
- ✅ Multi-stage Dockerfile (optimize image size)
- ✅ Alpine Linux base (minimal dependencies)
- ✅ Health checks & graceful shutdown
- ✅ Environment variable support

**Image Layers:**
```dockerfile
Stage 1: base         → Go 1.25 + dependencies
Stage 2: build        → Compile Go binary (5MB)
Stage 3: test-exec    → Run tests with coverage
Stage 4: final        → Alpine + binary (30MB)
```

---

### **Docker Compose**
```yaml
services:
  api:
    image: bookmark-api
    ports:
      - "8080:8080"
    environment:
      DB_HOST: postgres
      REDIS_URL: redis://redis:6379
  postgres:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: secret
  redis:
    image: redis:7
```

**Benefits:**
- ✅ Local development same as production
- ✅ All services in one command: `docker-compose up`
- ✅ Easy to onboard new developers

---

## 🧪 Testing & Quality

### **Test Coverage**
- ✅ Unit tests for services (mocked repositories)
- ✅ Integration tests for handlers
- ✅ Mock Redis for caching tests
- ✅ Test fixtures for consistent data

**Running Tests:**
```bash
go test ./...                           # Run all tests
go test -v ./...                        # Verbose output
go test -cover ./...                    # Coverage report
docker build --target=test .            # Run in Docker
```

---

### **Code Quality**
- ✅ Interface-based design (easy to mock)
- ✅ Dependency injection (decoupled components)
- ✅ Error handling (custom error types)
- ✅ Input validation (at handler level)
- ✅ SQL parameterized queries (prevent SQL injection)

---

## 📈 Performance Optimizations

### **1. Database Optimization**
- ✅ Indexes on frequently queried columns
- ✅ Pagination (avoid loading all records)
- ✅ GORM lazy loading control

**Example Query with Pagination:**
```go
// Bad: Load all bookmarks
var bookmarks []Bookmark
db.Where("user_id = ?", userID).Find(&bookmarks)

// Good: Load only 1 page (10 items)
var bookmarks []Bookmark
offset := (page - 1) * limit  // (1-1)*10 = 0
db.Where("user_id = ?", userID)
   .Offset(offset)
   .Limit(limit)
   .Find(&bookmarks)
```

### **2. Caching Strategy**
- ✅ Cache bookmark counts in Redis (reduce DB hits)
- ✅ Cache user sessions (faster authentication)
- ✅ TTL-based expiration (automatic cleanup)

### **3. Connection Pooling**
- ✅ GORM automatically pools DB connections
- ✅ Redis client connection pooling
- ✅ Configurable pool size

---

## 🔐 Security Best Practices

| Security Measure | Implementation |
|-----------------|-----------------|
| **Password Hashing** | bcrypt with salt rounds |
| **JWT Signing** | RS256 (RSA, not symmetric) |
| **Input Validation** | Struct tags, custom validators |
| **SQL Injection** | Parameterized queries via GORM |
| **XSS Protection** | JSON response format (no HTML) |
| **CORS** | Configurable origin whitelist |
| **Rate Limiting** | Redis-based throttling |
| **Soft Deletes** | Data preservation (never truly deleted) |

---

## 🔄 Deployment Flow

```
┌──────────────────┐
│   Git Push       │
└────────┬─────────┘
         │
┌────────▼──────────┐
│  GitHub Actions   │ (CI/CD)
│ • Run tests       │
│ • Build Docker    │
│ • Push to registry│
└────────┬──────────┘
         │
┌────────▼──────────┐
│  Docker Registry  │ (Docker Hub)
└────────┬──────────┘
         │
┌────────▼──────────┐
│  Deploy to Prod   │ (Kubernetes/VM)
│ • Pull image      │
│ • Start container │
│ • Health check    │
└──────────────────┘
```

---

## 📚 Project Structure

```
bookmark-management/
├── cmd/
│   ├── api/main.go              # Entry point
│   └── infrastructure/           # DI setup
│       ├── bootstrap.go          # Dependency injection
│       ├── database.go           # PostgreSQL init
│       ├── redis.go              # Redis init
│       └── config.go             # Config loading
│
├── internal/
│   ├── model/                    # Data entities
│   ├── repository/               # Database operations
│   ├── service/                  # Business logic
│   ├── handler/                  # HTTP handlers (Gin)
│   ├── api/                      # Routes & middleware
│   └── test/                     # Integration tests
│
├── pkg/
│   ├── sqldb/                    # DB utilities
│   ├── jwtutils/                 # JWT generation/validation
│   ├── redis/                    # Redis client
│   ├── logger/                   # Zerolog setup
│   └── utils/                    # Helper functions
│
├── migrations/                   # SQL migrations (.up.sql, .down.sql)
├── docs/                         # Swagger docs (auto-generated)
├── Dockerfile                    # Container image
├── docker-compose.yaml           # Local development setup
├── go.mod & go.sum               # Dependency management
└── README.md                     # Documentation
```

---

## 💡 Key Decisions & Trade-offs

### **Why PostgreSQL over MongoDB?**
- ✅ ACID transactions (data consistency)
- ✅ Strong relationships (foreign keys)
- ✅ Better for structured data (users, bookmarks)
- ❌ MongoDB better for unstructured/hierarchical data

### **Why JWT over Session-based?**
- ✅ Stateless (no session storage needed)
- ✅ Microservice-friendly (any server can validate)
- ✅ Mobile-friendly (API clients can store tokens)
- ❌ Sessions better for revocation (can't invalidate immediately)

### **Why Redis over in-memory cache?**
- ✅ Persistent across restarts
- ✅ Shareable between multiple servers
- ✅ Distributed systems support
- ❌ In-memory cache simpler for single server

### **Why Clean Architecture?**
- ✅ Easy to test (mock dependencies)
- ✅ Easy to modify (change one layer at a time)
- ✅ Scales with team growth
- ❌ Overkill for simple CRUD API

---

## 🎓 Learning Outcomes

### **Go Skills Developed**
- ✅ Interface-based design patterns
- ✅ Goroutines & concurrency basics
- ✅ Error handling strategies
- ✅ Testing with mocks and fixtures
- ✅ HTTP server development (Gin)
- ✅ Database integration (GORM, migrations)

### **Backend Architecture Skills**
- ✅ Clean Architecture implementation
- ✅ Dependency Injection pattern
- ✅ RESTful API design
- ✅ Authentication & authorization
- ✅ Database design & optimization
- ✅ Docker containerization

### **DevOps & Deployment**
- ✅ Multi-stage Docker builds
- ✅ Docker Compose for local development
- ✅ Database migrations strategy
- ✅ Configuration management (env vars)
- ✅ Health checks & monitoring

---

## 🚀 How to Present This in Interview

### **Opening Statement (30 seconds)**
> "I built a Bookmark Management REST API using Go and the Gin framework, following Clean Architecture principles. The project includes user authentication with JWT, PostgreSQL database with GORM ORM, Redis caching, and comprehensive test coverage. It's fully containerized with Docker and deployed using docker-compose for local development."

### **Architecture Explanation (1 minute)**
> "The project uses Clean Architecture with 4 distinct layers: handlers for HTTP requests, services for business logic, repositories for database access, and the database layer. This design ensures separation of concerns - if I need to change the database provider, I only modify the repository layer. I use dependency injection to pass dependencies, making it easy to mock repositories in unit tests."

### **Key Features (1 minute)**
> "The main features include user registration and login with JWT token generation using RS256 encryption, bookmark management with pagination, URL shortening, rate limiting using Redis, and CSV bulk import. All endpoints are protected except registration and login, using a JWT middleware that validates tokens and extracts user information."

### **Challenges & Solutions (1 minute)**
> "One challenge was managing database migrations in a team environment - I used golang-migrate with version-controlled SQL files for reproducible schema changes. Another was optimizing pagination queries - I added indexes on frequently queried columns and used GORM's offset/limit to avoid loading all records. For security, I implemented bcrypt password hashing and RSA-based JWT signing with separate private/public keys."

### **What I Would Do Differently**
- Implement refresh tokens (current tokens expire in 24h)
- Add request logging & distributed tracing
- Implement API versioning more explicitly
- Add comprehensive integration test suite
- Set up CI/CD pipeline (GitHub Actions)
- Add API rate limiting per user tier

---

## 🔗 Quick Reference

### **Common Commands**
```bash
# Setup
go mod tidy                              # Download dependencies
docker-compose up                        # Start services

# Running
go run cmd/api/main.go                   # Run server
go test ./...                            # Run all tests
docker build -t bookmark-api .           # Build image

# Testing
curl -X POST http://localhost:8080/users/login
curl -X GET http://localhost:8080/v1/bookmarks -H "Authorization: Bearer <token>"

# View API docs
http://localhost:8080/swagger/index.html
```

### **Important Files**
- `cmd/api/main.go` - Entry point
- `internal/api/api.go` - Routes & DI
- `internal/model/*.go` - Data models
- `internal/repository/*.go` - Database operations
- `internal/service/*.go` - Business logic
- `internal/handler/*.go` - HTTP handlers
- `pkg/jwtutils/*.go` - JWT logic
- `migrations/` - Database schema versions

---

## 📊 Tech Stack Summary Table

| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| **Language** | Go | 1.25.0 | Fast, concurrent backend |
| **Framework** | Gin Gonic | 1.12.0 | Web framework & routing |
| **Database** | PostgreSQL | Latest | Relational database |
| **ORM** | GORM | 1.31.1 | Database abstraction |
| **Auth** | JWT + RS256 | - | Stateless authentication |
| **Caching** | Redis | 9.18.0 | In-memory cache |
| **Testing** | Testify | 1.11.1 | Testing framework |
| **Logging** | Zerolog | 1.35.0 | Structured logging |
| **Config** | envconfig | 1.4.0 | Environment variables |
| **Docs** | Swagger | - | API documentation |
| **Container** | Docker | - | Containerization |

---

## ✅ Project Highlights

✨ **What Makes This Project Strong:**
1. ✅ Clean Architecture - professional code organization
2. ✅ Comprehensive Testing - unit & integration tests
3. ✅ Security First - bcrypt, JWT, input validation
4. ✅ Production Ready - Docker, migrations, logging
5. ✅ Well Documented - API docs, comments, guides
6. ✅ Scalable Design - can handle multiple users, pagination
7. ✅ Best Practices - SOLID principles, error handling
8. ✅ Modern Go - Using latest language features & packages

---

**Good luck with your interview! 🎉**
