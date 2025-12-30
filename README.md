# chuck
chuck norris api

## To Run
```sh
make dev
```

## API Endpoints

All endpoints return JSON.  
Base API path: `/api/v1`

---

## Health

### GET /health

Service health check.

**Auth:** Not required

**Example:**
```sh
curl http://localhost:8080/health
```

## Jokes

### GET /api/v1/jokes/random

Returns a random joke.

**Auth:** Not required

**Example:**
```sh
curl http://localhost:8080/api/v1/jokes/random
```

### GET /api/v1/jokes/search

Returns a random joke based on submitted query.

**Query Parameters:** query (string, required)
**Auth:** Required

**Example:**
```sh
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/v1/jokes/search?query=beard"
```

### GET /api/v1/jokes/personalized

Returns a random joke with submitted name for Chuck Norris.

**Query Parameters:** name (string, required)
**Auth:** Required

**Example:**
```sh
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8080/api/v1/jokes/personalized?name=Dave"
```

## Users and Auth

### POST /api/v1/users

Create a new user.

**Auth:** Not required

**Example:**
```sh
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

### POST /api/v1/auth/login

Authenticate a user and return an access token.

**Auth:** Not required

**Example:**
```sh
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```
