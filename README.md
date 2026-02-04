# Backend (be)

Go REST API backend using Clean Architecture.

## Tech Stack

- **Go 1.25**
- **PostgreSQL** (pgx/v5)
- **JWT** authentication (golang-jwt/v5)
- **godotenv** for environment config

## Project Structure

```
be/
├── cmd/api/          # Application entry point
├── internal/
│   ├── domain/       # Entities & interfaces
│   ├── repository/   # Data access layer
│   ├── service/      # Business logic
│   └── handler/      # HTTP handlers
├── go.mod
└── ARCHITECTURE.md   # Detailed architecture docs
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL

### Setup

1. Copy environment file:
```bash
cp .env.example .env
```

2. Update `.env` with your database credentials

3. Run the server:
```bash
go run cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /register | Register new user |
| POST | /login | User login |

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture documentation.
