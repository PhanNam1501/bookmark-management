# 🔐 JWT Authentication Guide

## Overview

This project uses **RSA-based JWT** for authentication:
- **Algorithm**: RS256 (RSA Signature with SHA-256)
- **Keys**: Private key (signing) + Public key (verification)
- **Expiration**: 24 hours

---

## Key Generation

### Generate Keys

```bash
# Generate 2048-bit RSA private key
openssl genrsa -out private.pem 2048

# Extract public key from private key
openssl rsa -in private.pem -pubout -out public.pem

# Verify (should show "RSA PUBLIC KEY")
cat public.pem
```

### File Locations
```
project-root/
├── private.pem  (Keep secret! Only server has this)
└── public.pem   (Public, can share)
```

---

## JWT Generation

### File: `pkg/jwtutils/generator.go`

```go
type JWTGenerator interface {
    GenerateToken(jwtContent jwt.MapClaims) (string, error)
}

func (g *jwtGenerator) GenerateToken(jwtContent jwt.MapClaims) (string, error) {
    // 1. Create token with RS256 algorithm
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
    
    // 2. Sign with private key
    tokenString, err := token.SignedString(g.privateKey)
    if err != nil {
        return "", err
    }
    
    return tokenString, nil
}
```

### Payload Format

```json
{
    "sub": "user-id-uuid",           // Subject (user ID)
    "iat": 1780419980,               // Issued at
    "exp": 1780506380                // Expiration
}
```

### JWT Token Structure

```
Header.Payload.Signature

Header:
{
    "alg": "RS256",
    "typ": "JWT"
}

Payload:
{
    "sub": "...",
    "iat": 1780419980,
    "exp": 1780506380
}

Signature: (signed with private key)
```

---

## JWT Validation

### File: `pkg/jwtutils/verifier.go`

```go
type JWTValidator interface {
    ValidateToken(tokenString string) (jwt.MapClaims, error)
}

func (j *jwtValidator) ValidateToken(tokenString string) (jwt.MapClaims, error) {
    // 1. Parse token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Return public key for verification
        return j.publicKey, nil
    })

    // 2. Check if token is valid
    if err != nil || !token.Valid {
        return nil, errInvalidToken
    }

    // 3. Extract claims
    return token.Claims.(jwt.MapClaims), nil
}
```

### Validation Steps
1. ✅ Parse JWT string
2. ✅ Verify signature with public key
3. ✅ Check expiration time
4. ✅ Return claims if valid

---

## Login Flow

### File: `internal/service/login.go`

```go
func (u *userService) Login(ctx context.Context, username, password string) (string, error) {
    // 1. Get user by username
    user, err := u.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return "", err
    }

    // 2. Verify password (bcrypt)
    passwordMatched := utils.VerifyPassword(user.Password, password)
    if !passwordMatched {
        return "", errors.New("wrong username or password")
    }

    // 3. Create JWT claims
    jwtContent := jwt.MapClaims{
        "sub": user.ID,                              // User ID
        "iat": time.Now().Unix(),                    // Now
        "exp": time.Now().Add(24 * time.Hour).Unix(), // 24h from now
    }

    // 4. Generate token (sign with private key)
    jwtString, err := u.jwtGen.GenerateToken(jwtContent)
    if err != nil {
        return "", err
    }

    return jwtString, nil
}
```

### Request/Response

```bash
# Request
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'

# Response
HTTP 200 OK
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## JWT Middleware

### File: `internal/api/middlewares/jwt.go`

```go
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // 1. Extract "Authorization: Bearer <token>" header
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required"})
            ctx.Abort()
            return
        }

        // 2. Parse "Bearer" + token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header"})
            ctx.Abort()
            return
        }

        tokenString := parts[1]

        // 3. Validate token (verify signature)
        claims, err := j.jwtValidator.ValidateToken(tokenString)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            ctx.Abort()
            return
        }

        // 4. Store claims in context (for handler to use)
        ctx.Set("claims", claims)

        // 5. Continue to next handler
        ctx.Next()
    }
}
```

### Middleware Usage

```go
// In internal/api/api.go
privateRoutes := a.app.Group("")
privateRoutes.Use(middlewares.NewJWTAuth(a.jwtValidator).JWTAuth())
privateRoutes.GET("/v1/bookmarks", handlers.BookmarkHandler.GetBookmarks)
```

---

## Using JWT in Handler

### Extract User ID from JWT

```go
func (h *handler) GetBookmarks(c *gin.Context) {
    // Get claims from context (set by middleware)
    claims, _ := c.Get("claims")
    mapClaims := claims.(jwt.MapClaims)
    
    // Extract user ID (sub = subject)
    userID := mapClaims["sub"].(string)
    
    // Now you have userID for querying user-specific data
    res, err := h.s.GetBookmarks(c, userID, limit, page)
    // ...
}
```

---

## Security Considerations

1. **Private Key Storage**
   - ✅ Keep private.pem secret
   - ✅ Never commit to version control
   - ✅ Use environment variables for key paths

2. **Token Expiration**
   - ✅ Set reasonable expiration (24h)
   - ✅ Implement refresh tokens (optional)
   - ✅ Validate expiration on every request

3. **Password Hashing**
   - ✅ Use bcrypt for password hashing
   - ✅ Never store plain text passwords
   - ✅ Verify hash on login

4. **HTTPS**
   - ✅ Always use HTTPS in production
   - ✅ Never send JWT over HTTP

5. **Token Validation**
   - ✅ Verify signature with public key
   - ✅ Check expiration time
   - ✅ Validate required claims

---

## Testing JWT

```bash
# 1. Register user
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123",
    "email": "john@example.com",
    "display_name": "John Doe"
  }'

# 2. Login (get token)
TOKEN=$(curl -s -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }' | jq -r '.token')

echo "Token: $TOKEN"

# 3. Use token to access protected endpoint
curl -X GET http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer $TOKEN"

# 4. Verify token rejection (missing header)
curl -X GET http://localhost:8080/v1/bookmarks
# Should return 401 Unauthorized

# 5. Verify token rejection (invalid token)
curl -X GET http://localhost:8080/v1/bookmarks \
  -H "Authorization: Bearer invalid-token"
# Should return 401 Unauthorized
```

---

## Debugging

### Decode JWT (inspect payload)

```bash
# Install jwt-cli (if available) or use online decoder
# https://jwt.io

# Or use jq to decode
echo "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..." | \
  cut -d. -f2 | \
  base64 -d | jq .
```

### Common Issues

| Issue | Solution |
|-------|----------|
| "Invalid token" | Check if private.pem/public.pem exist |
| "Missing auth header" | Add `Authorization: Bearer <token>` header |
| "Token expired" | Regenerate token (login again) |
| "Invalid signature" | Ensure public.pem matches private.pem |

---

## Related Files

- [Architecture Overview](ARCHITECTURE.md)
- [Database Setup](DATABASE.md)
- [API Endpoints](API.md)
