# 📊 Interview Diagrams & Visual References

## 0. 🏢 SYSTEM DESIGN - Full Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CLIENTS (Mobile/Web)                               │
│                                 │                                            │
│                                 ↓                                            │
│                        ┌────────────────┐                                   │
│                        │ Load Balancer  │ (Route requests)                   │
│                        │   (Nginx/HAProxy)                                  │
│                        └────────┬───────┘                                   │
│                                 │                                            │
│              ┌──────────────────┼──────────────────┐                        │
│              ↓                  ↓                  ↓                        │
│         ┌─────────┐        ┌─────────┐        ┌─────────┐                 │
│         │ API v1  │        │ API v2  │        │ API v3  │                 │
│         │(Go/Gin) │        │(Go/Gin) │        │(Go/Gin) │                 │
│         │Instance1│        │Instance2│        │Instance3│                 │
│         └────┬────┘        └────┬────┘        └────┬────┘                 │
│              │                  │                  │                        │
│              └──────────────────┼──────────────────┘                        │
│                                 │                                            │
│                    ┌────────────┴────────────┐                             │
│                    ↓                         ↓                             │
│           ┌──────────────────┐    ┌──────────────────┐                   │
│           │  PostgreSQL      │    │   Redis Cache    │                   │
│           │  (Main DB)       │    │   (Rate Limit,   │                   │
│           ├──────────────────┤    │   Session Cache) │                   │
│           │ - Users Table    │    │                  │                   │
│           │ - Bookmarks Tbl  │    │ Connection Pool: │                   │
│           │ - Indexes        │    │ 50-100 clients   │                   │
│           │ - Replication    │    └──────────────────┘                   │
│           └────────┬─────────┘                                            │
│                    │                                                       │
│                    ↓                                                       │
│           ┌──────────────────┐                                            │
│           │  PostgreSQL      │                                            │
│           │  (Read Replica)  │ ← Replication from main                    │
│           └──────────────────┘                                            │
│                                                                             │
│    ┌─────────────────────────────────────────────────────────────────┐   │
│    │         SUPPORTING SERVICES (Same instances)                    │   │
│    ├─────────────────────────────────────────────────────────────────┤   │
│    │                                                                 │   │
│    │  • JWT Token Generator    • Logger (Zerolog)                   │   │
│    │  • Password Hasher        • Error Tracker                      │   │
│    │  • Validator              • Health Check                       │   │
│    │  • Email Service (if any) • Metrics Collection                 │   │
│    │                                                                 │   │
│    └─────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                ┌───────────────────┼───────────────────┐
                ↓                   ↓                   ↓
         ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
         │   Monitoring │  │ Log Storage  │  │  Alerting    │
         │  (Prometheus)│  │   (ELK)      │  │  (PagerDuty) │
         └──────────────┘  └──────────────┘  └──────────────┘
```

---

### **System Components Explained:**

#### **1. Load Balancer**
```
Function: Distribute incoming requests to multiple API instances
Algorithm: Round-robin / Least connections
Benefits:
  ✓ High availability (if one instance fails, others handle requests)
  ✓ Horizontal scaling (add more instances as needed)
  ✓ Reduced single points of failure
```

#### **2. API Instances (Multiple)**
```
Why 3+ instances?
  ✓ Fault tolerance (if instance1 crashes, instance2 & 3 still serve)
  ✓ Parallel processing (handle more concurrent requests)
  ✓ Zero-downtime deployment (stop 1, update, restart)

Each instance contains:
  • Gin HTTP server
  • JWT validation middleware
  • All business logic
  • Connection pool to databases
```

#### **3. PostgreSQL (Primary + Replica)**
```
Primary (Write):
  → Receives all INSERT/UPDATE/DELETE
  → Replicates to read replica
  
Read Replica:
  → Serves SELECT queries (lighter load)
  → Read-only copy of data
  
Benefits:
  ✓ Scale reads (without affecting primary)
  ✓ Backup/disaster recovery
  ✓ High availability (failover if primary dies)
```

#### **4. Redis Cache**
```
Purpose: In-memory cache for frequently accessed data

Cache patterns:
  • User sessions: Fast authentication
  • Bookmark counts: Avoid expensive COUNT queries
  • Rate limit counters: Track API usage per user
  • JWT blacklist (optional): Revoke tokens

TTL Strategy:
  • Session: 24h (expires with JWT)
  • Counts: 1h (refresh periodically)
  • Rate limit: 1min (sliding window)

Benefits:
  ✓ Sub-millisecond response times
  ✓ Reduce database load
  ✓ Better user experience
```

---

### **Data Flow Examples:**

**Example 1: User Registration**
```
Client
  ↓ POST /users/register (username, password, email)
  ↓
Load Balancer
  ↓ Route to instance2 (round-robin)
  ↓
API Instance 2
  ├─ Validate input (handler layer)
  ├─ Hash password with bcrypt (service layer)
  ├─ Check if user exists (query cache first)
  ├─ INSERT into PostgreSQL (primary)
  ├─ Cache result in Redis
  ↓ Return HTTP 201 + user data
```

**Example 2: Get Bookmarks (Paginated)**
```
Client
  ↓ GET /v1/bookmarks?page=1&limit=10
  ↓ Header: Authorization: Bearer <JWT>
  ↓
Load Balancer
  ↓ Route to instance1
  ↓
API Instance 1
  ├─ Validate JWT (verify signature)
  ├─ Extract user_id from claims
  ├─ Check Redis cache for count
  │  ├─ HIT: Return cached count (fast!)
  │  └─ MISS: Query PostgreSQL
  ├─ Calculate pagination (offset, limit)
  ├─ Query PostgreSQL replica (read-only)
  ├─ Cache results for 1 minute
  ↓ Return HTTP 200 + bookmarks + pagination
```

**Example 3: Rate Limiting**
```
Client
  ↓ POST /v1/bookmarks (5th request in 1 minute)
  ↓
API Instance
  ├─ Extract user_id from JWT
  ├─ Check Redis: GET user:123:ratelimit
  │  → Current: 4 requests
  │  → Limit: 100/hour
  ├─ INCR counter
  ├─ TTL: 1 minute
  ├─ Check if exceeded
  │  ├─ YES: Return HTTP 429 Too Many Requests
  │  └─ NO: Process request
  ↓
```

---

## 1. 🏗️ Clean Architecture Layers

```
                      USER REQUEST
                            ↓
        ┌─────────────────────────────────────┐
        │   HTTP HANDLER LAYER (Gin)          │
        │  - Parse request (JSON)             │
        │  - Call service                     │
        │  - Return response (JSON)           │
        └──────────────┬──────────────────────┘
                       ↓
        ┌─────────────────────────────────────┐
        │   SERVICE LAYER (Business Logic)    │
        │  - Validation                       │
        │  - Authorization checks             │
        │  - Orchestration logic              │
        └──────────────┬──────────────────────┘
                       ↓
        ┌─────────────────────────────────────┐
        │   REPOSITORY LAYER (Data Access)    │
        │  - Query builder (GORM)             │
        │  - CRUD operations                  │
        │  - Database abstraction             │
        └──────────────┬──────────────────────┘
                       ↓
        ┌─────────────────────────────────────┐
        │   DATABASE LAYER (PostgreSQL)       │
        │  - Physical data storage            │
        │  - ACID transactions                │
        │  - Relationships (foreign keys)     │
        └─────────────────────────────────────┘
```

**Benefits:**
```
✓ If PostgreSQL → MySQL: Only change Repository layer
✓ Easy testing: Mock Repository layer
✓ Single responsibility: Each layer does one thing
```

---

## 2. 🔐 JWT Authentication Flow

```
┌──────────────────────────────────────────────────────┐
│ Step 1: LOGIN REQUEST                                │
└──────────────────┬───────────────────────────────────┘
                   │
                   POST /users/login
                   {
                     "username": "john",
                     "password": "secret123"
                   }
                   ↓
┌──────────────────────────────────────────────────────┐
│ Step 2: VERIFY PASSWORD (Service Layer)              │
│ - Query database: SELECT * FROM users WHERE username │
│ - Bcrypt.Compare(stored_hash, input_password)        │
│ - Result: ✓ Valid / ✗ Invalid                        │
└──────────────────┬───────────────────────────────────┘
                   ↓ (if valid)
┌──────────────────────────────────────────────────────┐
│ Step 3: CREATE JWT TOKEN                             │
│ Header:   {alg: "RS256", typ: "JWT"}                 │
│ Payload:  {sub: "user-id-123", iat: xxx, exp: xxx}   │
│ Signature: sign(header+payload, PRIVATE_KEY)         │
│ Result: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...      │
└──────────────────┬───────────────────────────────────┘
                   ↓
┌──────────────────────────────────────────────────────┐
│ Step 4: RETURN TOKEN TO CLIENT                       │
│ HTTP 200 OK                                          │
│ {                                                    │
│   "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..." │
│ }                                                    │
└──────────────────┬───────────────────────────────────┘
                   │
        ╔══════════════════════════════════════╗
        ║   CLIENT STORES TOKEN LOCALLY        ║
        ╚══════════════════════════════════════╝
                   │
                   ↓ (for next request)
┌──────────────────────────────────────────────────────┐
│ Step 5: PROTECTED REQUEST                            │
│ GET /v1/bookmarks?page=1&limit=10                    │
│ Header: Authorization: Bearer eyJhbGciOi...          │
└──────────────────┬───────────────────────────────────┘
                   ↓
┌──────────────────────────────────────────────────────┐
│ Step 6: JWT MIDDLEWARE VALIDATION                    │
│ 1. Extract token from header                         │
│ 2. Parse token (extract header.payload.signature)    │
│ 3. Verify signature: verify(payload, signature,      │
│                              PUBLIC_KEY) == true     │
│ 4. Check expiration: current_time < exp              │
│ 5. Extract claims: {sub: "user-id-123"}              │
└──────────────────┬───────────────────────────────────┘
                   ↓ (if valid)
┌──────────────────────────────────────────────────────┐
│ Step 7: PROCESS REQUEST                              │
│ - Handler receives user_id from JWT claims           │
│ - Query: SELECT * FROM bookmarks WHERE user_id = ?   │
│ - Return paginated results                           │
└──────────────────┬───────────────────────────────────┘
                   ↓
┌──────────────────────────────────────────────────────┐
│ Step 8: RETURN RESPONSE                              │
│ HTTP 200 OK                                          │
│ {                                                    │
│   "data": [bookmarks...],                            │
│   "pagination": {page: 1, limit: 10, total: 42}      │
│ }                                                    │
└──────────────────────────────────────────────────────┘
```

**Key Points:**
- 🔒 **Private key**: Only server has (signs tokens)
- 🔓 **Public key**: Can share (verifies tokens)
- ⏰ **Expiration**: 24 hours
- 🚫 **Cannot revoke**: Token valid until expiration (trade-off)

---

## 3. 💾 Database Schema Diagram

```
┌──────────────────────────────────┐
│          USERS TABLE             │
├──────────────────────────────────┤
│ id (PK, UUID)                    │◄─────────┐
│ username (UNIQUE)                │          │
│ email                            │          │
│ password_hash (bcrypt)           │          │
│ display_name                     │          │
│ created_at                       │          │
│ updated_at                       │          │
│ deleted_at (soft delete)         │          │
└──────────────────────────────────┘          │
         ▲                                     │ FK
         │                                     │
         │                                     │
         │                                     │
┌──────────────────────────────────┐          │
│       BOOKMARKS TABLE            │          │
├──────────────────────────────────┤          │
│ id (PK, UUID)                    │          │
│ user_id (FK) ──────────────────────────────┘
│ url (VARCHAR)                    │
│ description (TEXT)               │
│ code (UNIQUE, shortened URL id)  │
│ created_at                       │
│ updated_at                       │
│ deleted_at (soft delete)         │
└──────────────────────────────────┘

INDEXES:
├── users(username)       ← Fast user lookup
├── bookmarks(user_id)    ← Fast bookmark queries
└── bookmarks(code)       ← Fast URL redirect
```

**SQL Example:**
```sql
-- Get user's bookmarks
SELECT * FROM bookmarks
WHERE user_id = 'f1e27b78-97a5-4456-8163-6a83fade5dab'
  AND deleted_at IS NULL
LIMIT 10 OFFSET 0;
```

---

## 4. 🔄 Request Processing Pipeline

```
CLIENT REQUEST
      ↓
┌─────────────────────────────────────────┐
│ NETWORK LAYER (TCP/IP)                  │
│ - HTTP protocol handling                │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ GIN ROUTER                              │
│ Match: POST /v1/bookmarks → handler     │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ MIDDLEWARE CHAIN (executed first)       │
│ 1. CORS middleware                      │
│ 2. Logging middleware                   │
│ 3. JWT validation middleware            │
│    └─ Verify token signature ✓          │
│    └─ Extract claims (user_id) ✓        │
│    └─ Store in context                  │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ HANDLER LAYER                           │
│ Handler.CreateBookmark(c *gin.Context)  │
│ - Parse JSON body                       │
│ - Validate input                        │
│ - Extract user_id from JWT              │
│ - Call service                          │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ SERVICE LAYER                           │
│ Service.CreateBookmark(...)             │
│ - Additional validation                 │
│ - Business logic                        │
│ - Call repository                       │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ REPOSITORY LAYER (GORM)                 │
│ Repository.CreateBookmark(...)          │
│ - Build SQL: INSERT INTO bookmarks...   │
│ - Execute query                         │
│ - Return created record                 │
└────────────────┬────────────────────────┘
                 ↓
┌─────────────────────────────────────────┐
│ DATABASE LAYER (PostgreSQL)             │
│ - Physical storage                      │
│ - Write to disk                         │
│ - Return row                            │
└────────────────┬────────────────────────┘
                 ↓
        ← Response bubbles back ←
                 ↓
┌─────────────────────────────────────────┐
│ RESPONSE BUILDER                        │
│ - Format response {data: {...}}         │
│ - Set HTTP status 200                   │
│ - Serialize to JSON                     │
└────────────────┬────────────────────────┘
                 ↓
CLIENT RECEIVES: HTTP 200 + JSON body
```

---

## 5. 📦 Dependency Injection Pattern

```
Without DI (Tightly Coupled):
┌─────────────────┐
│ Handler         │
│  └─ creates     │
│    └─ Service   │
│      └─ creates │
│        └─ Repo  │
│          └─ creates
│            └─ DB Client
└─────────────────┘
PROBLEM: Hard to test, cannot mock


With DI (Loosely Coupled):
                    Bootstrap (main.go)
                          ↓
        ┌───────────────────┼───────────────────┐
        ↓                   ↓                   ↓
    ┌─────────┐         ┌─────────┐        ┌──────────┐
    │ Handler │◄────────│Service  │◄───────│Repository│
    └─────────┘(inject) └─────────┘(inject)└──────────┘
       ↑                    ↑                    ↑
       └────────────────────┴────────────────────┘
             All dependencies injected
             Easy to mock each component

BENEFIT: Handlers depend on interfaces, not concrete implementations
```

**Code Example:**
```go
// Bootstrap creates all dependencies
func Bootstrap() (*Engine, error) {
    db := sqldb.NewGORMClient(...)         // ← DB
    repo := repository.NewRepository(db)   // ← Inject DB
    svc := service.NewService(repo, ...)   // ← Inject Repo
    handler := handler.NewHandler(svc)     // ← Inject Service
    return NewEngine(handler), nil
}

// Now testing is easy
repo_mock := &MockRepository{}
svc := service.NewService(repo_mock)  // ← Inject mock
// Test service without touching real DB
```

---

## 6. 🔒 Password & JWT Security

```
╔════════════════════════════════════════════════════════════╗
║              PASSWORD HASHING (BCRYPT)                      ║
╠════════════════════════════════════════════════════════════╣
║                                                             ║
║  User Input: "mypassword123"                                ║
║       ↓                                                      ║
║  Bcrypt with salt (10 rounds)                               ║
║       ↓                                                      ║
║  Hash: "$2a$10$abcd1234efgh5678ijkl..."                     ║
║       ↓                                                      ║
║  Store in DB: users.password_hash                           ║
║                                                             ║
║  Login verification:                                        ║
║  bcrypt.Compare(stored_hash, input_password) → true/false  ║
║                                                             ║
║  WHY ONE-WAY?                                               ║
║  ✓ Cannot reverse hash back to password                     ║
║  ✓ If DB breached, hackers cannot use passwords            ║
║  ✓ Different hash for same password (salt)                 ║
║                                                             ║
╚════════════════════════════════════════════════════════════╝


╔════════════════════════════════════════════════════════════╗
║         JWT SIGNING (RS256 - ASYMMETRIC)                    ║
╠════════════════════════════════════════════════════════════╣
║                                                             ║
║  SIGNING (Server only):                                     ║
║  ┌─────────────────────────────────────┐                   ║
║  │ Claims: {sub: user_id, iat, exp}    │                   ║
║  └──────────────┬──────────────────────┘                   ║
║                 │                                           ║
║                 ├─ Encode header (algorithm, type)          ║
║                 │                                           ║
║                 ├─ Encode payload (claims)                  ║
║                 │                                           ║
║                 ├─ Sign with PRIVATE_KEY (in server only)  ║
║                 │                                           ║
║                 ↓                                           ║
║  Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...            ║
║                                                             ║
║  VERIFICATION (Middleware):                                 ║
║  ┌─────────────────────────────────────┐                   ║
║  │ Token from request header            │                   ║
║  └──────────────┬──────────────────────┘                   ║
║                 │                                           ║
║                 ├─ Extract header (algorithm)               ║
║                 │                                           ║
║                 ├─ Extract payload                          ║
║                 │                                           ║
║                 ├─ Extract signature                        ║
║                 │                                           ║
║                 ├─ Verify with PUBLIC_KEY (shared)          ║
║                 │                                           ║
║                 ├─ Check expiration: current < exp          ║
║                 │                                           ║
║                 ↓                                           ║
║  Result: ✓ Valid / ✗ Invalid (forged)                      ║
║                                                             ║
║  WHY ASYMMETRIC (RS256)?                                    ║
║  ✓ Private key signs (only server has it)                   ║
║  ✓ Public key verifies (can share, safe)                    ║
║  ✓ Token cannot be forged without private key              ║
║  ✓ Works across microservices                               ║
║                                                             ║
╚════════════════════════════════════════════════════════════╝
```

---

## 7. 🐳 Docker Multi-stage Build

```
┌──────────────────────────────────────┐
│ Stage 1: BASE                        │
│ FROM golang:1.25-alpine              │
│ + Install dependencies               │
│ Size: 500MB (temporary)              │
└──────────────────────────────────────┘
         ↓
┌──────────────────────────────────────┐
│ Stage 2: BUILD                       │
│ Compile Go binary                    │
│ CGO_ENABLED=0 go build               │
│ Binary size: 5MB                     │
└──────────────────────────────────────┘
         ↓
┌──────────────────────────────────────┐
│ Stage 3: TEST-EXEC (optional)        │
│ Run tests with coverage              │
│ Setup Redis for testing              │
│ Output: coverage.html                │
└──────────────────────────────────────┘
         ↓
┌──────────────────────────────────────┐
│ Stage 4: FINAL (Alpine)              │
│ FROM alpine:latest                   │
│ + Copy binary from stage 2 (5MB)     │
│ + Copy config files                  │
│ FINAL SIZE: ~30MB                    │
│                                      │
│ What's NOT included:                 │
│ ✗ Go compiler (saved 400MB)          │
│ ✗ Go source code (saved 50MB)        │
│ ✗ Test files (saved 10MB)            │
└──────────────────────────────────────┘
         ↓
    Docker Image Ready
    ~ 30MB (vs 500MB if single stage)
    + Contains only: binary, certs, migrations


BENEFIT: Smaller, faster to push/pull/deploy
```

---

## 8. 📊 API Response Format

```
SUCCESS RESPONSE:
┌────────────────────────────────────┐
│ HTTP 200 OK                        │
├────────────────────────────────────┤
│ {                                  │
│   "message": "Success message",    │
│   "data": {                        │
│     "id": "uuid",                  │
│     "user_id": "uuid",             │
│     "url": "...",                  │
│     "description": "...",          │
│     "created_at": "2026-06-03..."  │
│   },                               │
│   "pagination": {                  │
│     "page": 1,                     │
│     "limit": 10,                   │
│     "total": 42                    │
│   }                                │
│ }                                  │
└────────────────────────────────────┘

ERROR RESPONSE:
┌────────────────────────────────────┐
│ HTTP 401 Unauthorized              │
├────────────────────────────────────┤
│ {                                  │
│   "error": "Invalid token"         │
│ }                                  │
└────────────────────────────────────┘

VALIDATION ERROR:
┌────────────────────────────────────┐
│ HTTP 400 Bad Request               │
├────────────────────────────────────┤
│ {                                  │
│   "error": "Invalid input",        │
│   "details": {                     │
│     "url": "required",             │
│     "description": "max 500 chars" │
│   }                                │
│ }                                  │
└────────────────────────────────────┘
```

---

## 9. 🚀 Deployment Architecture

```
                    ┌─────────────────┐
                    │  Git Repository │
                    │   (GitHub)      │
                    └────────┬────────┘
                             │ Push
                             ↓
                    ┌─────────────────┐
                    │  CI/CD Pipeline │
                    │ (GitHub Actions)│
                    └────────┬────────┘
                             │
                    ┌────────┴────────┐
                    ↓                 ↓
              ✓ Tests          ✓ Build Docker
              ✓ Linters        ✓ Push to Registry
                             │
                             ↓
                    ┌─────────────────┐
                    │Docker Registry  │
                    │ (Docker Hub)    │
                    └────────┬────────┘
                             │ Pull
                             ↓
                    ┌─────────────────┐
                    │ Production      │
                    │ Kubernetes/VM   │
                    └─────────────────┘
```

**Local Development:**
```
docker-compose up
    │
    ├─ PostgreSQL container
    ├─ Redis container
    └─ API container

Then: go run cmd/api/main.go
```

---

## 10. 🧪 Testing Pyramid

```
                     △
                    ╱ ╲
                   ╱   ╲
                  ╱ E2E ╲          (1-2 tests)
                 ╱ Tests ╲         Integration tests
                ╱___────__ ╲       Testing full flow
               △
              ╱ ╲
             ╱   ╲
            ╱ Inte╲           (10-15 tests)
           ╱ grati╲          With mock DB
          ╱ on╲      ╲       Handler + Service
         ╱────____────╲
        △
       ╱ ╲
      ╱   ╲
     ╱ Unit╲               (50-100 tests)
    ╱ Tests ╲             Mocked dependencies
   ╱_________╲            Service layer only

EXAMPLE:
Unit Test → Mock Repository
Integration → Mock DB + Real Handler/Service
E2E Test → Real API + Real DB
```

---

## 11. 📈 Performance Optimization

```
PAGINATION:
┌────────────────────────────────────┐
│ Without pagination:                │
│ GET /bookmarks                     │
│ → Load ALL 1000 bookmarks          │
│ → Slow! (large response)           │
│ → Memory high! (1000 objects)      │
└────────────────────────────────────┘

┌────────────────────────────────────┐
│ With pagination:                   │
│ GET /bookmarks?page=1&limit=10     │
│ → Load only 10 bookmarks           │
│ → Fast! (small response)           │
│ → Memory low! (10 objects)         │
│ → User can navigate pages          │
└────────────────────────────────────┘

CACHING:
┌────────────────────────────────────┐
│ Count queries:                     │
│ SELECT COUNT(*) FROM bookmarks     │
│ WHERE user_id = ? ...              │
│                                    │
│ Problem: Slow with large tables    │
│ Solution: Cache in Redis (TTL 1h)  │
│ Result: Instant pagination info    │
└────────────────────────────────────┘

INDEXING:
┌────────────────────────────────────┐
│ Query: SELECT * FROM bookmarks     │
│        WHERE user_id = ?           │
│                                    │
│ Without index:                     │
│ → Full table scan (1000 rows)      │
│ → Slow!                            │
│                                    │
│ With index on user_id:             │
│ → Direct lookup (B-tree)           │
│ → Fast!                            │
└────────────────────────────────────┘
```

---

## 12. 🎯 Key Concepts Visualization

```
SOLID PRINCIPLES:
┌─────────────────────────────────────────────┐
│ S - Single Responsibility                   │
│   Each function/type has ONE reason to      │
│   change (Handler only handles HTTP,        │
│   Service only has business logic)          │
├─────────────────────────────────────────────┤
│ O - Open/Closed                             │
│   Open for extension (add new handlers)     │
│   Closed for modification (don't change     │
│   existing code)                            │
├─────────────────────────────────────────────┤
│ L - Liskov Substitution                     │
│   Can substitute implementations (if        │
│   interface says Repository, any            │
│   implementation should work)               │
├─────────────────────────────────────────────┤
│ I - Interface Segregation                   │
│   Small, focused interfaces (not one        │
│   big interface with everything)            │
├─────────────────────────────────────────────┤
│ D - Dependency Inversion                    │
│   High-level modules don't depend on        │
│   low-level modules, both depend on         │
│   abstractions (interfaces)                 │
└─────────────────────────────────────────────┘
```

---

---

## 13. 🚀 PRODUCTION DEPLOYMENT & SCALING ARCHITECTURE

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                        PRODUCTION ENVIRONMENT                                │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    CI/CD PIPELINE                                    │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                     │    │
│  │  GitHub Push                                                        │    │
│  │      ↓                                                              │    │
│  │  GitHub Actions                                                    │    │
│  │  ├─ Run go test ./...      (Unit tests)                            │    │
│  │  ├─ Run linters             (Code quality)                        │    │
│  │  ├─ Build Docker image      (Multi-stage)                         │    │
│  │  ├─ Push to registry        (Docker Hub)                          │    │
│  │  └─ Trigger deployment      (Kubernetes/VM)                       │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    DOCKER REGISTRY                                  │    │
│  │              (Docker Hub / Private Registry)                        │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                     │    │
│  │  Image: bookmark-api:v1.2.3                                        │    │
│  │  Size: ~30MB (optimized with multi-stage)                          │    │
│  │  Includes:                                                          │    │
│  │  ├─ Go binary (5MB)                                                │    │
│  │  ├─ Swagger docs                                                   │    │
│  │  ├─ Database migrations                                            │    │
│  │  └─ TLS certificates                                               │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │               KUBERNETES ORCHESTRATION                              │    │
│  │           (or Docker Compose for small scale)                       │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                     │    │
│  │  Pod/Container Orchestration:                                      │    │
│  │  ├─ Horizontal Pod Autoscaler (HPA)                               │    │
│  │  │  └─ Scale: 3-10 replicas based on CPU/Memory                  │    │
│  │  │                                                                 │    │
│  │  ├─ Service (Internal Load Balancer)                              │    │
│  │  │  └─ ClusterIP: 10.0.0.1 (internal DNS)                       │    │
│  │  │                                                                 │    │
│  │  ├─ Ingress (External Load Balancer)                              │    │
│  │  │  └─ External IP / Domain name                                 │    │
│  │  │     └─ TLS/HTTPS termination                                  │    │
│  │  │                                                                 │    │
│  │  └─ ConfigMap & Secrets                                           │    │
│  │     ├─ ConfigMap: Non-sensitive config                            │    │
│  │     └─ Secrets: API keys, DB credentials                          │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │              DATABASE CLUSTER (StatefulSet)                         │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                     │    │
│  │  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐          │    │
│  │  │   Postgres   │   │   Postgres   │   │   Postgres   │          │    │
│  │  │  Primary     │→→→│  Replica 1   │   │  Replica 2   │          │    │
│  │  │              │   │              │   │              │          │    │
│  │  │ Read/Write   │   │ Read-only    │   │ Read-only    │          │    │
│  │  │              │   │              │   │              │          │    │
│  │  └──────────────┘   └──────────────┘   └──────────────┘          │    │
│  │         ↑                   ↑                  ↑                   │    │
│  │  PersistentVolume (Storage)                                      │    │
│  │  (AWS EBS, Google Persistent Disk, etc)                          │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │              CACHING CLUSTER (Redis)                                │    │
│  ├─────────────────────────────────────────────────────────────────────┤    │
│  │                                                                     │    │
│  │  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐          │    │
│  │  │   Redis      │   │   Redis      │   │   Redis      │          │    │
│  │  │   Master     │→→→│   Slave 1    │   │   Slave 2    │          │    │
│  │  │              │   │              │   │              │          │    │
│  │  │  Sentinel for failover                            │          │    │
│  │  │              │   │              │   │              │          │    │
│  │  └──────────────┘   └──────────────┘   └──────────────┘          │    │
│  │                                                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘

                                    │
                ┌───────────────────┼───────────────────┐
                ↓                   ↓                   ↓
         ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
         │  Monitoring  │  │ Log Storage  │  │  Alerting    │
         │(Prometheus)  │  │  (ELK Stack) │  │ (PagerDuty)  │
         ├──────────────┤  ├──────────────┤  ├──────────────┤
         │ • CPU/Mem    │  │ • Centralize │  │ • Slack      │
         │ • Requests   │  │ • Search     │  │ • Email      │
         │ • Latency    │  │ • Visualize  │  │ • SMS        │
         │ • Errors     │  │   (Kibana)   │  │ • Webhook    │
         └──────────────┘  └──────────────┘  └──────────────┘
```

---

### **Scaling Strategies:**

#### **Horizontal Scaling (Add More Instances)**
```
CURRENT STATE:          AFTER SCALING:
┌─────────────┐        ┌─────────────┐
│ Load Bal.   │        │ Load Bal.   │
└──────┬──────┘        └──────┬──────┘
       │                      │
  ┌────┴────┐            ┌────┴────┬────────┐
  ↓         ↓            ↓    ↓    ↓        ↓
┌───┐  ┌───┐          ┌───┐┌───┐┌───┐   ┌───┐
│1  │  │2 │          │1 ││2 ││3 │   │4 │  (Add more instances)
└───┘  └───┘          └───┘└───┘└───┘   └───┘

✓ Handle more concurrent users
✓ Improve fault tolerance
✓ Auto-scale based on CPU/Memory
```

#### **Vertical Scaling (Bigger Machines)**
```
BEFORE:                 AFTER:
┌─────────────┐        ┌──────────────┐
│ 2 CPU cores │        │ 8 CPU cores  │
│ 4 GB RAM    │  →     │ 16 GB RAM    │
│ 50 conn/s   │        │ 500 conn/s   │
└─────────────┘        └──────────────┘

✓ Simpler (fewer instances to manage)
✗ Single point of failure
✗ Expensive (diminishing returns)
```

#### **Database Scaling**

**Read Replicas (for SELECT queries):**
```
App → Load Balancer → PostgreSQL Replica 1 (read)
                   → PostgreSQL Replica 2 (read)
                   → PostgreSQL Primary (write)
```

**Sharding (for massive scale):**
```
App → User ID hash → Shard 1 (users 1-1M)
                  → Shard 2 (users 1M-2M)
                  → Shard 3 (users 2M-3M)
```

---

### **Monitoring & Observability:**

```
METRICS (Prometheus):
├─ Request rate:     requests/second
├─ Latency (p50/p99):milliseconds
├─ Error rate:       5xx / 4xx errors
├─ CPU usage:        % of capacity
├─ Memory usage:     GB used
└─ DB connection:    Active pools

LOGS (ELK):
├─ Request logs:     Who, what, when
├─ Error logs:       Stack traces
├─ Audit logs:       API access
└─ Debug logs:       Development troubleshooting

TRACING (Jaeger/Datadog):
├─ Trace user request across services
├─ Identify bottlenecks
└─ Measure latency per component
```

---

### **Deployment Strategies:**

**Blue-Green Deployment:**
```
Blue (Current)      →  Green (New Version)
Production traffic  →  All traffic switched
                    →  Instant rollback if issue
```

**Canary Deployment:**
```
Version 1 (Old):    95% traffic
Version 2 (New):     5% traffic
                ↓ (if stable)
Version 1:           0% traffic
Version 2:         100% traffic
                ↓ (if issue, easy rollback)
```

**Rolling Deployment:**
```
Stop instance 1 → Update → Start
Stop instance 2 → Update → Start
Stop instance 3 → Update → Start
(Always 2/3 running, zero downtime)
```

---

**Print these diagrams or practice drawing them! 🎨**

Good luck! 🚀
