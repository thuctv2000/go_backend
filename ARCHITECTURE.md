# Architecture Documentation

This project follows **Clean Architecture** combined with the **Standard Go Project Layout**.

## Project Structure

```
my_backend/
├── cmd/
│   └── api/
│       └── main.go          # Entry point. Wires everything together.
├── internal/
│   ├── domain/              # Core Business Logic (Pure Go)
│   │   └── user.go          # Entities (Structs) and Interfaces (Repository/Service definitions)
│   ├── repository/          # Data Access Layer (Adapters)
│   │   └── memory_repo.go   # Implementation of domain.UserRepository (DB, Memory, etc.)
│   ├── service/             # Business Logic Layer (Use Cases)
│   │   └── auth_service.go  # Implementation of domain.AuthService (Register, Login flow)
│   └── handler/             # Transport Layer (HTTP/gRPC)
│       └── auth_handler.go  # Handles HTTP requests, parses JSON, calls Service.
└── go.mod
```

## Data Flow

`HTTP Request` -> `Handler` -> `Service` -> `Repository` -> `Database`

1.  **Handler**: Receives request, validates JSON, calls Service.
2.  **Service**: implementing business rules (hashing passwords, generating tokens), calls Repository.
3.  **Repository**: implementing database queries, returns Domain Entities.
4.  **Domain**: Defines *what* the data looks like and *what* operations are possible (Interfaces).

## How to Implement a New Feature (e.g., "Post" feature)

To add a new feature (e.g., creating and reading blog posts), follow these steps:

### 1. Domain Layer (`internal/domain/post.go`)
Define the Entity and the Interfaces.
```go
package domain

type Post struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

type PostRepository interface {
    Create(ctx context.Context, post *Post) error
    GetAll(ctx context.Context) ([]*Post, error)
}

type PostService interface {
    CreatePost(ctx context.Context, title, content string) (*Post, error)
    ListPosts(ctx context.Context) ([]*Post, error)
}
```

### 2. Repository Layer (`internal/repository/post_repo.go`)
Implement `PostRepository`.
```go
type memoryPostRepository struct { ... }
func NewMemoryPostRepository() domain.PostRepository { ... }
func (r *memoryPostRepository) Create(...) error { ... }
```

### 3. Service Layer (`internal/service/post_service.go`)
Implement `PostService`.
```go
type postService struct { repo domain.PostRepository }
func NewPostService(repo domain.PostRepository) domain.PostService { ... }
func (s *postService) CreatePost(...) (*domain.Post, error) { ... }
```

### 4. Handler Layer (`internal/handler/post_handler.go`)
Implement HTTP handlers.
```go
type PostHandler struct { service domain.PostService }
func NewPostHandler(s domain.PostService) *PostHandler { ... }
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) { ... }
```

### 5. Wiring (`cmd/api/main.go`)
Connect everything in `main.go`.
```go
postRepo := repository.NewMemoryPostRepository()
postService := service.NewPostService(postRepo)
postHandler := handler.NewPostHandler(postService)

mux.HandleFunc("POST /posts", postHandler.Create)
```
