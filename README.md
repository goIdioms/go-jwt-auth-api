# go-jwt-auth-api

A RESTful API for user authentication and authorization using JWT, built with Go, Fiber, MongoDB, and Redis.

## Features

- User registration and login
- JWT-based authentication (access & refresh tokens)
- Role-based access control
- Secure password hashing
- Token refresh endpoint
- Logout functionality
- User profile and users listing
- Swagger (OpenAPI) documentation
- Docker & Docker Compose support
- Redis integration for caching/session management

## Tech Stack

- **Go** (Golang)
- **Fiber** web framework
- **MongoDB** for data storage
- **Redis** for caching/session
- **JWT** for authentication
- **Swagger** for API docs
- **Docker** for containerization

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### Running Locally

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/go-jwt-auth-api.git
    cd go-jwt-auth-api
    ```

2. Copy `.env.example` to `.env` and set your environment variables.

3. Start the services:
    ```bash
    make build
    make run
    ```

4. The API will be available at `http://localhost:8080/api`
5. Swagger docs: `http://localhost:8080/swagger/index.html`

### API Endpoints

- `POST /api/auth/signup` — Register a new user
- `POST /api/auth/signin` — Login and get tokens
- `POST /api/auth/refresh` — Refresh tokens
- `POST /api/auth/logout` — Logout

- `GET /api/auth/me` — Get current user
- `GET /api/users` — List users(admin, moderator only)
