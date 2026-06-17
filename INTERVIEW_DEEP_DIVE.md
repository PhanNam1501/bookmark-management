# 🎯 INTERVIEW DEEP DIVE - WHAT / WHY / WHAT BROKE

Hướng dẫn chuẩn bị chi tiết theo 3 lớp: **WHAT** (cái gì), **WHY** (tại sao), **WHAT BROKE** (cái gì hỏng).

---

## 1. 🏗️ CLEAN ARCHITECTURE

### **WHAT - Bạn xây dựng gì?**

Tôi tổ chức dự án thành 4 layers riêng biệt:
- **Handler Layer**: Xử lý HTTP requests, parse JSON, return responses
- **Service Layer**: Chứa business logic, validation, orchestration
- **Repository Layer**: Database operations, CRUD
- **Database Layer**: Physical storage (PostgreSQL)

**Kết quả đo được:**
- ✅ 80%+ code coverage (via unit tests)
- ✅ Thêm tính năng mới mà không ảnh hưởng code cũ
- ✅ Dễ dàng swap PostgreSQL → MySQL (chỉ change Repository layer)
- ✅ Test cases không phải touch database

---

### **WHY - Tại sao bạn chọn nó?**

**Trade-off:**
```
CLEAN ARCHITECTURE:
✓ Dễ test (mock dependencies)
✓ Dễ sửa đổi (change one layer)
✓ Scalable (add features easily)
✓ Clear responsibilities
✗ Boilerplate code (more files)
✗ Overkill cho simple CRUD API
✗ Learning curve (phải hiểu patterns)

ALTERNATIVES:
1. MVC (Model-View-Controller):
   ✓ Simpler structure
   ✗ Controllers thường quá fat (300+ lines)
   ✗ Business logic spreads everywhere
   → Good for small projects, not scalable

2. Monolithic (No structure):
   ✓ Quick to start
   ✗ Hard to test
   ✗ Hard to maintain (grows into spaghetti code)
   → Okay for first prototype, not production
```

**Lý do chọn Clean Architecture:**
> "Tôi biết project này sẽ grow thêm features (URL shortening, CSV import, rate limiting). Nếu dùng sloppy structure, business logic sẽ scatter everywhere. Clean Architecture buộc tôi think in layers - nếu thêm feature, tôi biết chính xác modify cái gì (Repository? Service? Handler?). Disadvantage là boilerplate, nhưng cái đó worth it cho maintainability."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: Quá nhiều layers cho đơn giản features**
```
Example: Simple "Get User by ID"

Naive approach:
GET /users/123 → Handler → Service → Repository → DB
(4 layers for 1 query)

Một junior dev comment: "Tại sao phải qua 4 layers?
Có thể gọi DB trực tiếp trong handler mà!"

Vấn đề: Nếu logic phức tạp, sẽ phải refactor nhiều lần
```

**Cách fix:**
> "Đúng, cho simple queries có thể gọi trực tiếp. Nhưng trong practice, 'simple' lại grow complexity (thêm caching, authorization checks, logging). Nên tôi accept 'boilerplate' này từ đầu - nó là investment cho future."

**Vấn đề 2: Interface overuse**
```
Ban đầu tôi tạo interface cho tất cả:

type Handler interface {
    CreateBookmark(c *gin.Context)
    GetBookmarks(c *gin.Context)
    DeleteBookmark(c *gin.Context)
}

Problem: Gin không implement interface này
Solution: Interface chỉ dùng cho dependency injection
         Không cần interface cho tất cả (YAGNI - You Ain't Gonna Need It)
```

**Cách fix:**
> "Tôi học được: Only create interfaces khi bạn actually need to mock hoặc swap implementations. Nếu không, interface chỉ adds complexity. Tôi giữ lại interface cho Repository (vì cần mock trong tests) nhưng bỏ Handler interface."

---

## 2. 🔐 JWT AUTHENTICATION (RS256)

### **WHAT - Bạn xây dựng gì?**

Tôi implement stateless authentication system:
- User login → bcrypt password → generate JWT → return token
- Protected requests → verify JWT signature → extract user_id → query user-specific data
- JWT expires sau 24h

**Architecture:**
```
pkg/jwtutils/generator.go  → Sign token with RS256 (private key)
pkg/jwtutils/verifier.go   → Verify token (public key)
internal/api/middlewares/jwt.go → Extract & validate token

Kết quả:
✓ No session storage needed (stateless)
✓ Works across multiple servers
✓ Mobile/SPA friendly
✓ Can't forge token without private key
```

---

### **WHY - Tại sao bạn chọn JWT + RS256?**

**Trade-off:**
```
JWT (STATELESS):
✓ Scalable (no session storage)
✓ Microservice-friendly (any server can validate)
✓ Mobile-friendly (store token locally)
✓ Self-contained (claims inside token)
✗ Can't revoke immediately (token valid until expiration)
✗ Bigger request size (token in every request)
✗ Token theft = full access (no session refresh)

SESSION-BASED:
✓ Can revoke immediately
✓ More secure (session stored server-side)
✓ Smaller request size
✗ Requires session storage (Redis/DB)
✗ Not scalable across servers (sticky sessions)
✗ Hard to implement for mobile
```

**RS256 (ASYMMETRIC) vs HS256 (SYMMETRIC):**
```
RS256 (Asymmetric):
✓ Public key can be shared (verify without exposing secret)
✓ Better for distributed systems
✗ Slower (RSA operations)

HS256 (Symmetric):
✓ Faster (HMAC operations)
✓ Simpler setup
✗ Secret key shared everywhere (security risk)
✗ Can't revoke public key
```

**Lý do chọn JWT + RS256:**
> "Dự án này có thể grow sang microservices (Auth service, Bookmark service). JWT stateless nên any service có thể validate token without hitting Auth service. RS256 cho phép public key sharing - Auth service public public.pem, other services verify locally. Plus, stateless = dễ scale horizontally (thêm server instances không cần share session store)."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: Token không thể revoke**
```
Scenario: User logout, nhưng token vẫn valid cho đến expiration

Hướng dẫn xác thực: "Bạn logout lúc 2:00 PM, nhưng token hết hạn lúc 2:00 AM (24h later).
Nếu token leak, hacker có quyền truy cập 24 giờ!"

Impact: High-security apps không thể chỉ dùng JWT (cần thêm blacklist)
```

**Cách fix:**
```
Option 1: Token blacklist (Redis)
  POST /logout → Add token to Redis blacklist
  Middleware: Check if token in blacklist
  TTL: Same as token expiration (24h)
  Cost: Extra Redis lookup per request

Option 2: Short-lived tokens + Refresh tokens
  Access token: 15 minutes
  Refresh token: 7 days
  User logout: Blacklist refresh token
  Cost: Extra logic, more complex

Option 3: Accept risk (cho non-critical apps)
  "This is password manager app, not banking.
  Accept risk of 24h token validity."
```

**Cách tôi handle:**
> "Tôi implement Redis blacklist cho logout. Not ideal (extra Redis call), but acceptable trade-off. Nếu cần ultra-high-security, tôi implement refresh tokens (15min access + 7day refresh)."

**Vấn đề 2: Key rotation complications**
```
Scenario: Cần rotate private key (security best practice every 6 months)

Problem:
  1. Old tokens signed with old key sẽ invalid
  2. Clients logged out suddenly
  3. No way to transition smoothly

In production:
  Solution: Maintain multiple public keys
  Server signs new tokens with new private key
  Verify using: old_public_key || new_public_key
  Graceful transition period (6 months)
```

**Vấn đề 3: Timing attack vulnerability**
```
Ban đầu verify JWT như này:

if token == stored_token {
    return true
}

Problem: String comparison time depends on number of matching chars
  Comparing "abc" vs "xyz" → fast
  Comparing "abcdefgh" vs "abcdefgh" → slow
Hacker can measure time → guess token byte-by-byte

Cách fix: Dùng time-constant comparison
  crypto/subtle.ConstantTimeCompare(a, b)
  Same time regardless of matching characters
```

---

## 3. 💾 DATABASE: PostgreSQL + GORM

### **WHAT - Bạn xây dựng gì?**

Tôi setup PostgreSQL với:
- **2 tables**: Users, Bookmarks (with soft deletes)
- **Foreign keys**: Bookmarks.user_id → Users.id
- **Indexes**: 
  - users(username) - for login
  - bookmarks(user_id) - for user's bookmarks
  - bookmarks(code) - for URL redirect
- **Migrations**: golang-migrate (version-controlled schema changes)

**Kết quả:**
```
✓ Query bookmarks by user_id: ~1ms (with index)
✓ Login: ~2ms (username index + bcrypt compare)
✓ Pagination: avoid loading all records
✓ Soft deletes: recover deleted data
```

---

### **WHY - Tại sao PostgreSQL + GORM?**

**Trade-off:**
```
POSTGRESQL:
✓ ACID transactions (data consistency)
✓ Foreign keys (relationships)
✓ Powerful query language (complex queries)
✓ Mature ecosystem
✓ Open source, free
✗ Vertical scaling limits (need replicas for read heavy)
✗ Slower writes than MongoDB

MONGODB:
✓ Horizontal scaling (built-in sharding)
✓ Flexible schema (easy schema changes)
✓ Fast writes
✗ No transactions (until recent versions)
✗ No foreign keys (application-level joins)
✗ Eventual consistency (not ACID)

MYSQL:
✓ Faster than PostgreSQL (for writes)
✗ Less powerful features
✗ Less mature

GORM vs Raw SQL:
✓ Type-safe queries
✓ Auto migration support
✓ Less boilerplate
✗ Performance overhead (small)
✗ Learning curve
```

**Lý do chọn PostgreSQL + GORM:**
> "Dữ liệu users/bookmarks có relationships (1-to-many). ACID transactions important cho consistency (user register + create default bookmark must both succeed or both fail). PostgreSQL perfect cho structured data. GORM cho phép auto-generate migrations (track schema changes in git) và type-safe queries (prevent SQL injection automatically)."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: N+1 Query Problem**
```
Code:

for _, bookmark := range bookmarks {
    user := db.Query("SELECT * FROM users WHERE id = ?", bookmark.UserID)
    // Process user
}

Problem: 
  SELECT bookmarks (1 query)
  SELECT user for bookmark 1 (query 2)
  SELECT user for bookmark 2 (query 3)
  ...
  SELECT user for bookmark 100 (query 101)
  
Total: 101 queries instead of 2!

Symptom: Thêm bookmarks → exponentially slow queries
```

**Cách fix:**
```
Option 1: Eager Loading (Preload)
  db.Preload("User").Find(&bookmarks)
  → 2 queries: 1 for bookmarks, 1 for all users (JOIN)

Option 2: Select specific columns
  db.Select("id", "user_id", "url").Find(&bookmarks)
  → Reduce data transfer

Option 3: Use joins
  db.Joins("JOIN users ON bookmarks.user_id = users.id").Find(&bookmarks)
```

**Cách tôi handle:**
> "Tôi use GORM's Preload() - đơn giản nhất. Học được: always think about query efficiency. Thêm test case với 1000 bookmarks để catch N+1 trước production."

**Vấn đề 2: Missing Indexes = Slow Queries**
```
Ban đầu không có index trên user_id:

SELECT * FROM bookmarks WHERE user_id = 'user-123'

Without index: Full table scan (1M rows)
  Time: ~500ms ❌

With index:
  Time: ~1ms ✓

Impact: User click bookmarks → wait 500ms
```

**Cách fix:**
```sql
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
```

**Cách tôi handle:**
> "Tôi identify hot queries (Login = WHERE username, Get bookmarks = WHERE user_id) và add indexes immediately. Lesson: measure first (use EXPLAIN ANALYZE) trước optimize. Nếu query dưới 10ms, không cần optimize."

**Vấn đề 3: Migration Order Matters**
```
Scenario:

Migration 1: CREATE TABLE users
Migration 2: CREATE TABLE bookmarks
Migration 3: ALTER TABLE bookmarks ADD FOREIGN KEY (user_id) REFERENCES users(id)

If developer apply migrations wrong order:
  Apply migration 3 first → fails (users table doesn't exist)

Problem: No dependency tracking
```

**Cách fix:**
```
Migration naming convention: 001_create_users.up.sql
                           002_create_bookmarks.up.sql
                           003_add_fk.up.sql

golang-migrate enforces sequential order ✓
```

---

## 4. ⚡ REDIS CACHING

### **WHAT - Bạn xây dựng gì?**

Tôi implement caching layer cho:
- **Rate limiting**: Track requests per user (INCR, TTL)
- **Session cache**: Cache user auth info (key: "user:{id}", TTL: 24h)
- **Bookmark count**: Cache count for pagination (TTL: 1h)

**Architecture:**
```
Request → Check Redis (1ms) → HIT: return cached value
                            → MISS: query DB (50ms) + cache result

Kết quả:
✓ Pagination info: 1ms (vs 50ms from DB)
✓ Rate limit check: 1ms
✓ Reduce DB load: 80% hit rate on counts
```

---

### **WHY - Tại sao bạn chọn Redis caching?**

**Trade-off:**
```
REDIS CACHING:
✓ Very fast (sub-millisecond)
✓ Simple key-value interface
✓ Distributed (shareable between instances)
✓ Atomic operations (INCR, DECR)
✗ Extra infrastructure (another service to maintain)
✗ Data loss if Redis crashes (no persistence)
✗ Stale data (cache hit = maybe old value)

IN-MEMORY CACHE (Go map):
✓ No extra infrastructure
✓ No network latency
✗ Not shared between instances
✗ Data lost on restart
✗ No TTL support

DATABASE CACHING (PostgreSQL):
✓ Persistent
✗ Much slower (500x slower than Redis)
✗ Not suitable for frequently-accessed data

NO CACHING:
✓ Always fresh data
✗ Slow (wait for DB)
✗ High DB load
✗ Bad user experience
```

**Lý do chọn Redis:**
> "Bookmark count queries lặp lại nhiều (mỗi pagination request). Database COUNT(*) expensive trên large tables. Redis cho phép cache counts với TTL 1h - sau 1h auto-expire, tôi chấp nhận data cũ vài giây. Distributed - works with horizontal scaling (multiple API instances)."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: Cache Invalidation Nightmare**
```
Scenario: User create bookmark → update count cache

Code:
1. INSERT INTO bookmarks (...)
2. Redis.INCR("count:{user_id}")

Problem: Nếu step 2 fail, count cache sai
  DB: +1 bookmarks
  Redis: không change
  Result: count off by 1

Multiply with many bookmarks → count completely wrong
```

**Cách fix:**
```
Option 1: Invalidate cache on write
  INSERT → Redis.DEL("count:{user_id}")
  Next request → query DB fresh → cache new value
  
Option 2: Consistent caching (transactions)
  Use Redis WATCH for optimistic locking
  Use PostgreSQL + Redis transaction together
  
Option 3: TTL-based (eventual consistency)
  Set SHORT TTL (5min) instead of 1h
  Accept stale data up to 5min
  Trade: More cache misses (slower) vs accuracy
```

**Cách tôi handle:**
> "Tôi use simple invalidation: CREATE bookmark → Redis.DEL(count_key) → next pagination request queries fresh count. Trade-off: More DB hits pero guaranteed accuracy. Nếu traffic high, tôi set TTL để auto-expire."

**Vấn đề 2: Cache Stampede**
```
Scenario: Cache expire + 100 concurrent requests

Timeline:
  00:00 - Cache hit count: 42
  01:00 - Cache expire
  01:00.001 - User1 request → Cache MISS → Query DB
  01:00.002 - User2 request → Cache MISS → Query DB
  01:00.003 - User3 request → Cache MISS → Query DB
  ...
  01:00.100 - User100 request → Query DB

Result: 100 concurrent DB queries (spike!)
vs normal: 1 query (single cache MISS)
```

**Cách fix:**
```
Option 1: Probabilistic early expiration
  TTL: 1h ± random(0-5min)
  Prevent all keys expiring at same time

Option 2: Refresh cache in background
  Cron job: Refresh hot keys every 50min
  vs wait for expiration

Option 3: Set stale flag
  Keep value until expired + 1h
  If TTL < 0, return stale value + refresh in background
```

**Cách tôi handle:**
> "Tôi monitor Redis hit rate. Nếu stampede detect, tôi implement probabilistic TTL (TTL = 1h + random 0-10min) hoặc background refresh job."

**Vấn đề 3: Redis Crash = Cascading Failure**
```
Scenario: Redis down

User request:
  1. Redis unreachable (timeout 1s)
  2. Fallback to DB query (slow)
  3. Other users also fallback
  4. DB overwhelmed (thundering herd)
  5. API crashes

Result: Single point of failure!
```

**Cách fix:**
```
Option 1: Graceful degradation
  If Redis timeout → query DB (slower but works)
  Log error (monitor)
  Don't fail user request

Option 2: Redis cluster
  Multiple Redis instances (master + replicas)
  Sentinel for failover
  
Option 3: Circuit breaker
  Track Redis errors
  After N errors → skip Redis (fast fail)
  After recovery → resume Redis
```

**Cách tôi handle:**
> "Tôi wrap Redis calls với timeout (500ms) + fallback to DB. Nếu Redis slow, better return data slow than crash completely. Plus monitor Redis health - alert if down > 5min."

---

## 5. 🧪 TESTING STRATEGY

### **WHAT - Bạn xây dựng gì?**

Tôi implement 3 levels of testing:
```
UNIT TESTS (Services):
- Mock Repository
- Test business logic isolated
- 50+ tests

INTEGRATION TESTS (Handlers + Services):
- Mock database
- Test request → handler → service → mock repo → response
- Use Testify for assertions
- 20+ tests

END-TO-END TESTS (with real DB):
- Docker compose start real services
- Test full flow
- 5+ tests (slow, only critical paths)
```

**Kết quả:**
```
✓ 80%+ code coverage
✓ Catch regressions early
✓ Confidence to refactor
✓ Documentation via tests
```

---

### **WHY - Tại sao bạn chọn test strategy này?**

**Trade-off:**
```
UNIT TESTS (Isolated, Fast):
✓ Fast (run 100 tests in 1s)
✓ Easy to write
✓ Catch logic bugs quickly
✗ Don't catch integration bugs
✗ Mock everything

INTEGRATION TESTS (More realistic):
✓ Catch integration bugs
✓ More realistic than unit
✗ Slower
✗ Harder to debug

END-TO-END TESTS (Real systems):
✓ Most realistic
✓ Catch real bugs
✗ Very slow
✗ Flaky (external dependencies)

NO TESTS:
✓ Fast to write code
✗ No safety net
✗ Scary to refactor
✗ Production bugs
```

**Pyramid approach (many unit, few E2E):**
```
     /\
    /E2E\        (5 tests, slow)
   /──────\
  / Inte-  \   (20 tests, medium)
 /──────────\
/  Unit     \  (50+ tests, fast)
```

**Lý do chọn strategy này:**
> "Tôi prioritize unit tests (cover service logic) vì họ run fast. Integration tests catch handler bugs. E2E tests catch real-world issues. Pyramid shape = most ROI per test. Nếu test failure, unit test fail first (easy to debug), không phải waiting for slow E2E."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: Mock Repository Lies**
```
Scenario: Mock Repository return success for DELETE

MockRepo.DeleteBookmark(ctx, id) → nil

Test pass! 🎉

Reality: PostgreSQL cascade delete rule fail
         Constraints violated
         Actual delete fail in production

Problem: Mocks hide real database bugs
```

**Cách fix:**
```
Option 1: Integration test with real DB
  Use testcontainers (Docker container for Postgres in test)
  Real constraints + triggers validate
  Cost: Tests slower (1-2s each)

Option 2: Add more unit test scenarios
  Test cascade delete rules in unit tests
  Mock database errors (timeout, constraint violation)
  
Option 3: E2E test critical paths
  Test DELETE via API → real DB
```

**Cách tôi handle:**
> "Tôi add integration tests cho data operations (DELETE, CREATE with FKs). Trade-off: Tests slower nhưng catch real bugs. Plus E2E test for critical flows (register + create bookmark)."

**Vấn đề 2: Test Data Inconsistency**
```
Test 1:
  Create user "john"
  Create bookmark for john
  ✓ PASS

Test 2:
  Create user "john" (same name)
  ✗ FAIL (unique constraint)

Problem: Test order matters! 
Test depends on previous test state
This is BAD - tests should be independent
```

**Cách fix:**
```
Solution: Test fixtures + cleanup

func TestCreateBookmark(t *testing.T) {
    // Setup: Fresh DB state
    db.ExecSQL("TRUNCATE bookmarks, users")
    user := createTestUser("john")
    
    // Test
    err := service.CreateBookmark(...)
    assert.NoError(t, err)
    
    // Cleanup
    db.ExecSQL("TRUNCATE bookmarks, users")
}
```

**Cách tôi handle:**
> "Tôi implement fixtures (helper functions) để setup test data consistent. Plus run tests with `-shuffle` flag (randomize order) để catch ordering bugs early."

**Vấn đề 3: Flaky Tests (Intermittent Failures)**
```
Test sometimes pass, sometimes fail (no code change)

Common causes:
1. Race conditions (concurrent access)
2. Timing issues (sleep 1s not enough)
3. External dependencies (Redis timeout)
4. Random data generation
```

**Cách fix:**
```
Option 1: Use testify assertions
  assert.Eventually(t, func() bool {
      return isConditionTrue()
  }, 5*time.Second, 100*time.Millisecond)
  
Option 2: Mock time for timing tests
  Use mocktime library (control time in tests)
  
Option 3: Eliminate external dependencies
  Use miniredis (mock Redis) instead of real Redis
  Use testcontainers (isolated DB per test)
```

**Cách tôi handle:**
> "Tôi use miniredis (mock Redis) để tests không depend on real Redis. Use testify Eventually for async assertions. Plus run tests 10x (`for i in {1..10}; do go test ./...; done`) to catch flakiness."

---

## 6. 🐳 DOCKER & DEPLOYMENT

### **WHAT - Bạn xây dựng gì?**

Tôi create:
- **Multi-stage Dockerfile**: base → build → test → final (30MB image)
- **docker-compose.yaml**: PostgreSQL + Redis + API locally
- **Health checks**: Endpoint `/health` để verify container running
- **Environment variables**: Configurable per environment (dev/staging/prod)

**Kết quả:**
```
✓ Deploy anywhere Docker runs (consistency)
✓ New devs: docker-compose up → ready (5 min)
✓ No "works on my machine" problems
✓ Small image (30MB vs 500MB non-optimized)
```

---

### **WHY - Tại sao bạn chọn Docker?**

**Trade-off:**
```
DOCKER:
✓ Consistency (same image dev/staging/prod)
✓ Isolation (no conflicts with other apps)
✓ Easy scaling (start N containers)
✓ Easy CI/CD (build once, deploy everywhere)
✗ Performance overhead (small, ~5%)
✗ Learning curve
✗ Storage overhead (images + layers)

BARE METAL / VMs:
✓ Full control
✓ Better performance
✗ Setup complexity
✗ Hard to scale
✗ Dependency hell (different versions)

KUBERNETES:
✓ Advanced orchestration
✗ Overkill for small projects
✗ Complex (learning curve)
```

**Lý do chọn Docker:**
> "Docker guarantee consistency: image built in CI works same in staging & production. Multi-stage build optimize image size (final 30MB, not 500MB). docker-compose cho local development - new team members setup in 5 minutes instead of hours configuring services."

---

### **WHAT BROKE - Cái gì hỏng?**

**Vấn đề 1: Docker Image Size Explosion**
```
Ban đầu single-stage Dockerfile:

FROM golang:1.25
COPY . .
RUN go build -o api cmd/api/main.go

EXPOSE 8080
CMD ["./api"]

Result: Image size = 500MB 😱
  Go compiler: 400MB
  Source code: 50MB
  Binary: 5MB
  
Problem: Push/pull image slow (1 min)
         Deploy slow
         Storage expensive
```

**Cách fix:**
```
Multi-stage Dockerfile:

Stage 1 (builder): Compile binary (in 500MB image)
Stage 2 (final): FROM scratch / alpine
                 COPY binary from stage 1 (5MB)
                 
Result: Final image = 30MB ✓ (6x smaller!)
        Push/pull: 5s (vs 60s before)
        Deploy: faster
```

**Cách tôi handle:**
> "Tôi implement multi-stage build từ đầu. Stage 1 compile, stage 2 copy binary to scratch/alpine. Lesson: optimize early, don't wait until production complains."

**Vấn đề 2: Database Connection Leaks**
```
Scenario: Container running, but connection pool exhausted

GORM default:
  MaxOpenConns: unlimited
  ConnMaxLifetime: unlimited

Code:
  for i := 0; i < 1000; i++ {
      db.Query(...) // open connection
  }
  // Missing close → connections leak

Result: After 1 hour, 1000 open connections
        → DB connection pool full
        → New requests hang (timeout)
```

**Cách fix:**
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(50)      // Max 50 concurrent
sqlDB.SetMaxIdleConns(10)      // Keep 10 idle
sqlDB.SetConnMaxLifetime(time.Hour) // Recycle old connections

// OR use context properly:
db.WithContext(ctx).Query(...) // Auto cleanup
```

**Cách tôi handle:**
> "Tôi set connection pool limits explicitly (MaxOpenConns=50). Test with load testing (50 concurrent users) to catch connection leaks. Plus monitoring: if connection count keeps growing = leak detected."

**Vấn đề 3: Graceful Shutdown Timeout**
```
Scenario: Docker stop container

DEFAULT behavior:
  1. SIGTERM sent
  2. Container waits 10 seconds (--stop-timeout=10s)
  3. SIGKILL force kill

Problem: Request in-flight during shutdown
  User sends POST /bookmarks
  Server gets SIGTERM
  Response never sent (connection killed)
  
User sees: connection reset error ❌
```

**Cách fix:**
```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           handler,
    ReadTimeout:       15 * time.Second,
    WriteTimeout:      15 * time.Second,
    IdleTimeout:       60 * time.Second,
    MaxHeaderBytes:    1 << 20,
    GracefulStopDelay: 30 * time.Second, // Gracefully stop
}

// In main:
go func() {
    <-sigChan // SIGTERM
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    server.Shutdown(ctx)
}()
```

**Cách tôi handle:**
> "Tôi implement graceful shutdown: on SIGTERM, wait 30 seconds for in-flight requests to finish, then force kill. Plus set Docker stop-timeout=35s. Result: no dropped requests during deployments."

---

## 7. 🎯 BONUS: Lessons Learned & What I'd Do Differently

### **Git Commits & Code Review**

**WHAT BROKE:**
```
Commit message: "fix bugs"

6 months later, looking at git blame:
  "Why did this code change? What bug?"
  → Can't tell from commit message
  → Waste time investigating

Also: Commit cách quá lớn (500 lines)
      Hard to review
      Hard to revert if problematic
```

**Cách tôi fix:**
```
Commit guidelines:
1. One logical change per commit
2. Commit message: "What" + "Why"
   - Bad: "fix pagination"
   - Good: "fix pagination offset calculation for page 1 edge case - was returning page 2 results"

3. Small commits (< 200 lines)
   → Easy to review
   → Easy to revert if issue
   → Easy to understand history
```

---

### **Error Handling**

**WHAT BROKE:**
```
Code:
  user, _ := repo.GetUserByID(ctx, id)  // ignore error
  
Problem: If user not found = nil
         Next line: user.Name → panic (nil pointer)
         Service crashed

Horror: This error only happen when user deleted mid-request
        Rare case, not caught in testing
```

**Cách tôi fix:**
```go
// Check errors explicitly
user, err := repo.GetUserByID(ctx, id)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return errors.New("user not found")
    }
    return fmt.Errorf("failed to get user: %w", err) // wrap with context
}

// Custom error types
var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized")

// Handle in handler
if errors.Is(err, ErrUserNotFound) {
    c.JSON(404, gin.H{"error": "user not found"})
} else if errors.Is(err, ErrUnauthorized) {
    c.JSON(401, gin.H{"error": "unauthorized"})
} else {
    c.JSON(500, gin.H{"error": "internal server error"})
}
```

---

### **Performance Optimization**

**WHAT BROKE:**
```
User: "API slow, takes 2 seconds to fetch bookmarks"

Investigation:
  - Handler: 10ms
  - Service: 15ms
  - Repository: 500ms ← BOTTLENECK!
  
Cause: SELECT * FROM bookmarks WHERE user_id = ?
       No index on user_id
       Full table scan (1M rows) = slow
```

**Cách tôi fix:**
```
1. Measure first:
   db.WithContext(ctx).Explain("EXPLAIN ANALYZE", &benchmarks).Find(&bookmarks)
   → See execution plan

2. Add index:
   CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id)

3. Verify improvement:
   500ms → 5ms ✓ (100x faster!)

Key lesson: Profile before optimize
           Don't guess, measure
           Optimize hot paths only
```

---

## 📋 Summary: How to Prepare for Different Questions

### **"Tell me about your authentication system"**
- **WHAT**: JWT + RS256, 24h expiration, stateless
- **WHY**: Scalable (no session store), works across microservices, public key can be shared safely
- **WHAT BROKE**: Can't revoke immediately → add Redis blacklist for logout

### **"How would you optimize slow queries?"**
- **WHAT**: Identified N+1 problem with bookmarks → users
- **WHY**: Each bookmark loaded user separately (100 bookmarks = 101 queries)
- **WHAT BROKE**: Tested with 1000 bookmarks, query timeout → added GORM Preload()

### **"Why Clean Architecture?"**
- **WHAT**: 4 layers (handler → service → repo → db)
- **WHY**: Easy to test (mock each layer), easy to change (swap PostgreSQL)
- **WHAT BROKE**: Initially created interfaces for everything → realized only mock repositories needed

### **"How do you ensure code quality?"**
- **WHAT**: 80%+ code coverage with unit + integration tests
- **WHY**: Pyramid approach (many unit, few E2E) balances speed & confidence
- **WHAT BROKE**: Mocks hid database bugs → added integration tests with real DB constraints

---

**This is the level of detail that impresses interviewers! 🚀**

Use this framework for ANY project / ANY technology:
1. Explain WHAT you built (concrete, measurable)
2. Explain WHY you chose it (trade-offs, alternatives, business logic)
3. Explain WHAT BROKE (real problems, how you debugged, lessons learned)

Good luck! 💪
