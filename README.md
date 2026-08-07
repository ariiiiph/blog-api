# Blog API

A production-style RESTful Blog API built with **Go** and **SQLite**.

This project provides a complete backend system for a blogging platform with user authentication, authorization, post management, social interactions, image uploads, search, pagination, and Swagger API documentation.

The project follows a layered architecture with separated responsibilities between routes, services, repositories, and database layers.

---

# Features

## Authentication & Authorization

* User registration and login
* JWT authentication
* Refresh token mechanism
* Protected routes
* Role-based authorization
* Admin role support
* Password hashing

## Blog Management

* Create posts
* Edit posts
* Delete posts
* Categories management
* Tags management
* Comments system
* Like system
* Image upload for posts

## Advanced Features

* Search functionality
* Pagination
* Swagger API documentation
* Middleware-based request handling
* Database migrations
* Docker containerization
* Docker Compose support
* Persistent storage using Docker volumes

---

# Tech Stack

* **Backend:** Go
* **Database:** SQLite
* **Authentication:** JWT + Refresh Tokens
* **API Documentation:** Swagger
* **Configuration:** cleanenv
* **Containerization:** Docker + Docker Compose
* **Architecture:** Layered Architecture
* **API Testing:** Postman

---

# Project Architecture

The project follows a layered backend architecture:

```
HTTP Request

     ↓

Routes

     ↓

Services

     ↓

Repositories

     ↓

Database
```

### Responsibilities

### Routes

* Handle HTTP requests and responses

### Services

* Contain business logic

### Repositories

* Handle database operations

### Models

* Define application entities

### Middlewares

* Handle authentication, CORS, and request processing

---

# Project Structure

```
backend/

├── cmd/
│   └── api/
│       └── main.go
│
├── config/
│   └── dev.env
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── db/
│   │   ├── db.go
│   │   └── migrations/
│   │
│   ├── middlewares/
│   ├── models/
│   ├── repository/
│   ├── routes/
│   ├── service/
│   └── utils/
│
├── sqlite/
├── uploads/
│   └── posts/
│
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

---

# Database

The project uses **SQLite** with version-controlled SQL migrations.

Implemented database features:

* Users table
* Posts table
* Categories
* Tags
* Post tags relationship
* Comments
* Likes
* Refresh tokens
* User roles
* Post images

Migration files are located in:

```
internal/db/migrations/
```

---

# Configuration

The application supports two configuration methods.

## Local Development

For running the project directly with Go:

```
config/dev.env
```

Example:

```env
ENV=dev

DB_PATH=sqlite/dev
DB_NAME=api.db

HTTP_ADDRESS=localhost:8080

JWT_KEY=your_secret_key
```

Run:

```bash
go run ./cmd/api -config config/dev.env
```

---

## Docker Environment

When running with Docker Compose, environment variables are injected into the container.

Create a `.env` file based on:

```
.env.example
```

Example:

```env
ENV=dev

DB_PATH=sqlite/dev
DB_NAME=api.db

HTTP_ADDRESS=0.0.0.0:8080

JWT_KEY=your_secret_key
```

---

# Running The Project

## Requirements

* Go 1.26+
* Docker
* Docker Compose
* Git

---

# Run Locally

Clone repository:

```bash
git clone https://github.com/ariiiiph/blog-api.git
```

Move into project:

```bash
cd blog-api
```

Install dependencies:

```bash
go mod download
```

Run:

```bash
go run ./cmd/api -config config/dev.env
```

The API will start at:

```
http://localhost:8080
```

---

# Run With Docker

Build and start the application:

```bash
docker compose up --build
```

The API will be available at:

```
http://localhost:8080
```

Swagger:

```
http://localhost:8080/swagger/
```

Stop containers:

```bash
docker compose down
```

---

# Docker Storage

The application uses Docker volumes to persist data.

Volumes:

```
blog-data
```

Stores:

```
SQLite database
```

---

```
blog-uploads
```

Stores:

```
Uploaded post images
```

This allows data to survive container recreation.

---

# Authentication

The API uses JWT authentication.

For protected routes, include the access token:

```http
Authorization: Bearer YOUR_ACCESS_TOKEN
```

Refresh tokens are used to generate new access tokens after expiration.

---

# Swagger Documentation

Swagger documentation is included in:

```
docs/
```

Available files:

```
swagger.json
swagger.yaml
docs.go
```

Swagger provides:

* Available API endpoints
* Request formats
* Response formats
* API testing interface

---

# API Features Overview

| Feature               | Status |
| --------------------- | ------ |
| User Registration     | ✅      |
| Login                 | ✅      |
| JWT Authentication    | ✅      |
| Refresh Tokens        | ✅      |
| Admin Role            | ✅      |
| Posts CRUD            | ✅      |
| Categories            | ✅      |
| Tags                  | ✅      |
| Comments              | ✅      |
| Likes                 | ✅      |
| Search                | ✅      |
| Pagination            | ✅      |
| Image Upload          | ✅      |
| Swagger Documentation | ✅      |
| Docker Support        | ✅      |
| Docker Compose        | ✅      |
| Persistent Volumes    | ✅      |

---

# Example Request

## Create Post

```json
{
    "title": "My First Blog Post",
    "content": "Hello from my Go Blog API",
    "category_id": 1,
    "tags": [
        "golang",
        "backend"
    ]
}
```

---

# File Upload

The API supports image uploads for posts.

Uploaded images are stored locally:

```
uploads/posts/
```

When using Docker, uploads are persisted using Docker volumes.

---

# Future Improvements

* Add automated unit and integration tests
* Add PostgreSQL support
* Add Redis caching
* Add CI/CD pipeline
* Deploy to AWS
* Add monitoring and logging
* Improve security hardening
* Add rate limiting

---

# License

This project is created for learning purposes and as a backend engineering portfolio project.
