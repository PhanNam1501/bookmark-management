# 📡 API Endpoints Documentation

## Base URL
```
http://localhost:8080
```

---

## Authentication

All endpoints marked with 🔒 require JWT token in header:
```
Authorization: Bearer <your_token>
```

---

## User Endpoints

### Register User
```
POST /users/register
```

**Request:**
```json
{
  "username": "john",
  "password": "password123",
  "email": "john@example.com",
  "display_name": "John Doe"
}
```

**Response (200 OK):**
```json
{
  "id": "f1e27b78-97a5-4456-8163-6a83fade5dab",
  "username": "john",
  "email": "john@example.com",
  "display_name": "John Doe",
  "created_at": "2026-06-03T00:00:00Z"
}
```

---

### Login 🔓
```
POST /users/login
```

**Request:**
```json
{
  "username": "john",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Invalid username or password"
}
```

---

### Get Current User 🔒
```
GET /v1/self/info
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "f1e27b78-97a5-4456-8163-6a83fade5dab",
  "username": "john",
  "email": "john@example.com",
  "display_name": "John Doe",
  "created_at": "2026-06-03T00:00:00Z"
}
```

---

## Bookmark Endpoints

### Create Bookmark 🔒
```
POST /v1/bookmarks
```

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "url": "https://www.google.com",
  "description": "Google Search"
}
```

**Response (200 OK):**
```json
{
  "message": "Bookmark created successfully",
  "data": {
    "id": "dd679abf-c9fb-48a4-bbd2-c58f0950253",
    "user_id": "f1e27b78-97a5-4456-8163-6a83fade5dab",
    "url": "https://www.google.com",
    "description": "Google Search",
    "code": "f94128e1032209442b3e1bf6e8ef60",
    "created_at": "2026-06-03T00:00:00Z"
  }
}
```

---

### Get Bookmarks (Paginated) 🔒
```
GET /v1/bookmarks?page=1&limit=10
```

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (required): Page number (starts from 1)
- `limit` (required): Items per page

**Response (200 OK):**
```json
{
  "message": "Bookmarks retrieved successfully",
  "data": [
    {
      "id": "dd679abf-c9fb-48a4-bbd2-c58f0950253",
      "user_id": "f1e27b78-97a5-4456-8163-6a83fade5dab",
      "url": "https://www.google.com",
      "description": "Google Search",
      "code": "f94128e1032209442b3e1bf6e8ef60",
      "created_at": "2026-06-03T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 42
  }
}
```

---

## Testing with cURL

### 1. Register User
```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123",
    "email": "john@example.com",
    "display_name": "John Doe"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'
```

Save the token from response:
```bash
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 3. Get Current User
```bash
curl -X GET http://localhost:8080/v1/self/info \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Create Bookmark
```bash
curl -X POST http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.google.com",
    "description": "Google Search"
  }'
```

### 5. Get Bookmarks
```bash
curl -X GET "http://localhost:8080/v1/bookmarks?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Using Postman

### Import Collection

**User Register**
- Method: `POST`
- URL: `http://localhost:8080/users/register`
- Body (JSON):
  ```json
  {
    "username": "john",
    "password": "password123",
    "email": "john@example.com",
    "display_name": "John Doe"
  }
  ```

**User Login**
- Method: `POST`
- URL: `http://localhost:8080/users/login`
- Body (JSON):
  ```json
  {
    "username": "john",
    "password": "password123"
  }
  ```

**Get Bookmarks** (Protected)
- Method: `GET`
- URL: `http://localhost:8080/v1/bookmarks?page=1&limit=10`
- Headers:
  - Key: `Authorization`
  - Value: `Bearer <token_from_login>`

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid query parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "unauthorized"
}
```

### 500 Internal Server Error
```json
{
  "error": "failed to get bookmarks",
  "details": "connection refused"
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK - Success |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing/invalid token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

---

## Pagination

### Request
```
GET /v1/bookmarks?page=2&limit=20
```

### Response
```json
{
  "data": [...],
  "pagination": {
    "page": 2,
    "limit": 20,
    "total": 100
  }
}
```

**Calculation:**
- `offset = (page - 1) * limit`
- `offset = (2 - 1) * 20 = 20`
- Returns items 21-40

---

## Swagger Documentation

Interactive API documentation available at:
```
http://localhost:8080/swagger/index.html
```

---

## Common Workflows

### Workflow 1: Register and Get Bookmarks

```bash
# 1. Register
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "securepass",
    "email": "alice@example.com",
    "display_name": "Alice"
  }'

# 2. Login
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "securepass"}' > login.json

# Extract token
TOKEN=$(cat login.json | jq -r '.data.token')

# 3. Get bookmarks
curl -X GET "http://localhost:8080/v1/bookmarks?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### Workflow 2: Create Multiple Bookmarks

```bash
TOKEN="your_token_here"

# Create bookmark 1
curl -X POST http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com", "description": "GitHub"}'

# Create bookmark 2
curl -X POST http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://stackoverflow.com", "description": "Stack Overflow"}'

# Create bookmark 3
curl -X POST http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com", "description": "Google"}'

# List all bookmarks
curl -X GET "http://localhost:8080/v1/bookmarks?page=1&limit=100" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Related Files

- [Architecture Overview](ARCHITECTURE.md)
- [JWT Guide](JWT.md)
- [Database Setup](DATABASE.md)
