# ⚡ Interview Cheat Sheet - Quick Reference

## 🎯 Elevator Pitch (30 seconds)
> "I built a Bookmark Management REST API using Go + Gin framework. It follows Clean Architecture with JWT authentication (RS256), PostgreSQL + GORM for database, Redis for caching, and is fully containerized with Docker. Production-ready with tests and Swagger docs."

---

## 🏗️ Architecture (4 Layers)
```
Handler (Gin) → Service (Logic) → Repository (DB) → Database
```

**Why Clean Architecture?**
- Separation of concerns
- Easy to test (mock each layer)
- Easy to modify (change one layer without affecting others)

---

## 🛠️ Core Tech Stack

### **Required (Must Know)**
| Tech | Why |
|------|-----|
| **Go 1.25** | Fast, concurrent, simple syntax |
| **Gin Gonic** | Lightweight web framework |
| **PostgreSQL** | ACID transactions, relationships |
| **GORM** | Type-safe ORM, migrations |
| **JWT + RS256** | Stateless auth with asymmetric crypto |
| **Bcrypt** | Password hashing (one-way) |
| **Docker** | Containerization, reproducibility |

### **Additional (Good to Know)**
| Tech | Purpose |
|------|---------|
| Redis | Rate limiting, caching |
| Zerolog | Structured logging |
| Testify | Unit testing |
| Swagger | API documentation |

---

## 🔐 Authentication Flow

```
POST /users/login
{username, password}
    ↓
Service: Verify password (bcrypt)
    ↓
Generate JWT Claims: {sub: user_id, iat, exp}
    ↓
Sign with private.pem (RS256)
    ↓
Return token string
```

**Protected Request:**
```
GET /v1/bookmarks
Header: Authorization: Bearer <token>
    ↓
Middleware: Verify signature with public.pem
    ↓
Extract user_id from claims
    ↓
Query bookmarks for that user
```

---

## 📊 Database Schema (Simple)

```
Users:
├── id (UUID)
├── username (UNIQUE)
├── email
├── password_hash (bcrypt)
└── created_at, updated_at, deleted_at

Bookmarks:
├── id (UUID)
├── user_id (FK)
├── url
├── description
├── code (short identifier)
└── created_at, updated_at, deleted_at
```

**Indexes:**
- ✅ username (fast login)
- ✅ user_id on bookmarks (fast bookmark queries)
- ✅ code (fast URL redirect)

---

## 🚀 Key Features

**1. User Management**
- Register user
- Login (get JWT token)
- Get current user info

**2. Bookmark CRUD**
- Create bookmark
- List bookmarks (paginated)
- Delete bookmark (soft delete)

**3. Additional**
- URL shortening
- Rate limiting (Redis)
- CSV import
- Password management

---

## 🧪 Testing Strategy

```
Unit Tests:
→ Test services (mock repositories)
→ Test handlers (mock services)

Integration Tests:
→ Test with real/mock database
→ Use Testify for assertions
```

**Run Tests:**
```bash
go test ./...                    # All tests
go test -cover ./...             # With coverage
docker build --target=test .     # In Docker
```

---

## 🐳 Docker Multi-stage Build

```dockerfile
Stage 1: base         → Go + dependencies
Stage 2: build        → Compile binary
Stage 3: test-exec    → Run tests + coverage
Stage 4: final        → Alpine + binary only
```

**Benefits:**
- Small final image (30MB)
- Tests run in CI/CD
- Production-ready

---

## 📈 Performance Optimizations

| Problem | Solution |
|---------|----------|
| Load all bookmarks | Use pagination (offset/limit) |
| Slow login | Cache user sessions in Redis |
| Count queries | Cache in Redis with TTL |
| No indexes | Index on user_id, username, code |
| N+1 queries | GORM eager loading |

---

## 🔒 Security Measures

✅ **Passwords:** Hashed with bcrypt (one-way)  
✅ **JWT:** Signed with RS256 (cannot forge without private key)  
✅ **User Isolation:** Each user only sees their own data  
✅ **Input Validation:** Struct tags + custom validators  
✅ **SQL Injection:** Parameterized queries (GORM)  
✅ **Soft Deletes:** Data preserved, never truly deleted  
✅ **Rate Limiting:** Redis-based throttling  

---

## 💡 Key Decisions

| Question | Answer | Why |
|----------|--------|-----|
| Why PostgreSQL? | ACID + relationships | Better than MongoDB for structured data |
| Why JWT? | Stateless | Microservice-friendly, works with mobile |
| Why Redis? | Distributed cache | Persistent, shareable between servers |
| Why Clean Architecture? | Testability | Easy to mock, modify, scale |

---

## 🎓 Technical Concepts Demonstrated

- ✅ Interface-based design (Go)
- ✅ Dependency Injection (decoupling)
- ✅ Clean Architecture (layered design)
- ✅ RESTful API design
- ✅ Database design & normalization
- ✅ JWT authentication
- ✅ Password hashing (bcrypt)
- ✅ Pagination implementation
- ✅ Docker containerization
- ✅ Unit testing & mocking
- ✅ Goroutines & concurrency
- ✅ Error handling strategies

---

## ❓ Common Interview Questions

### Q1: "Why did you choose Go?"
> "Go is fast, has excellent built-in concurrency with goroutines, and compiles to a single binary making deployment easy. It's perfect for backend services and microservices."

### Q2: "How does authentication work?"
> "Users login with username/password. We verify the password using bcrypt, then create a JWT token with RS256 signature using a private key. For protected endpoints, we verify the signature with the public key and extract the user_id from claims."

### Q3: "What's your architecture?"
> "Clean Architecture with 4 layers: Handler (HTTP), Service (business logic), Repository (database), Database. This allows easy testing and modification - if I change the database, only Repository layer changes."

### Q4: "How do you handle pagination?"
> "Calculate offset = (page-1) * limit, use GORM with Offset() and Limit(). Also query total count separately. Returns current page, limit, and total for the frontend to render pagination controls."

### Q5: "What about security?"
> "Passwords are hashed with bcrypt (one-way). JWTs are signed with RS256 (asymmetric). Each user only sees their own data. We use parameterized queries to prevent SQL injection. Soft deletes preserve data."

### Q6: "How do you test?"
> "Unit tests for services using mocked repositories. Integration tests for handlers with mock database. Use Testify for assertions. Docker for isolated test environments with Redis."

### Q7: "What would you do differently?"
> "Implement refresh tokens. Add distributed tracing. Set up CI/CD pipeline. More comprehensive integration tests. API versioning. Implement different rate limits per user tier."

---

## 🗂️ File Structure Quick Map

```
cmd/api/main.go              ← Start here
  ↓
cmd/infrastructure/bootstrap.go    ← Dependency injection
  ↓
internal/api/api.go               ← Routes definition
  ↓
internal/handler/*.go             ← HTTP handlers (Gin)
  ↓
internal/service/*.go             ← Business logic
  ↓
internal/repository/*.go          ← Database operations
  ↓
internal/model/*.go               ← Data structures
```

---

## 📋 Before Interview Checklist

- [ ] Understand the 4-layer architecture
- [ ] Know JWT flow by heart
- [ ] Explain why each tech was chosen
- [ ] Prepare examples (login, create bookmark)
- [ ] Know database schema
- [ ] Explain testing strategy
- [ ] Have answers ready for common questions
- [ ] Practice elevator pitch (30 seconds)
- [ ] Can draw architecture diagram

---

## 🎯 Key Talking Points

1. **"Clean Architecture"** - explain separation of concerns
2. **"Dependency Injection"** - explain testability benefit
3. **"JWT with RS256"** - explain asymmetric cryptography
4. **"Bcrypt"** - explain one-way hashing
5. **"Pagination"** - explain offset-limit pattern
6. **"Docker Multi-stage"** - explain size optimization
7. **"Soft Deletes"** - explain data preservation
8. **"Indexes"** - explain query optimization

---

## ⏱️ Time Breakdown

| Topic | Time | What to Say |
|-------|------|------------|
| **Project Overview** | 30s | Name, tech, purpose |
| **Architecture** | 1m | 4 layers, why clean arch |
| **Authentication** | 1m | JWT + RS256 + bcrypt flow |
| **Database** | 1m | Schema, indexes, migrations |
| **Features** | 1m | Register, login, CRUD bookmarks |
| **Testing** | 30s | Unit + integration tests |
| **Deployment** | 30s | Docker, docker-compose |
| **Challenges** | 1m | Problems & solutions |
| **Learnings** | 1m | What you learned & gained |

**Total: ~8 minutes** (leaves room for questions)

---

**Pro Tips:**
- 🎤 Speak clearly, don't rush
- 📊 Use hand gestures for architecture
- 💬 Tell stories about problems you solved
- ⚡ Show enthusiasm for the tech
- 🤔 Ask clarifying questions if confused
- 😊 Be honest about what you don't know

Good luck! 🚀
