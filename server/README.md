# AP School - Backend

Go REST API for the AP School platform.

## Requirements

- Go 1.23+
- Docker
- [Goose](https://github.com/pressly/goose) (for migrations)

## Setup

### 1. Environment Variables

```bash
cp .env.example .env
# Edit .env with your values
```

### 2. Start Database

```bash
make docker-run
```

### 3. Run Migrations

```bash
make migrate-up
```

### 4. Seed Data (optional)

```bash
make seed
```

### 5. Start Server

```bash
# With hot reload
make watch

# Without hot reload
make run
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Start the server |
| `make watch` | Start with hot reload (requires air) |
| `make build` | Build binary |
| `make test` | Run tests |
| `make docker-run` | Start PostgreSQL container |
| `make docker-down` | Stop PostgreSQL container |
| `make migrate-up` | Apply migrations |
| `make migrate-down` | Revert last migration |
| `make migrate-status` | Check migration status |
| `make seed` | Load challenges into database |

## Project Structure

```
server/
├── cmd/
│   ├── api/
│   │   ├── main.go          # Entry point
│   │   └── routes.go        # Route definitions
│   └── seed/
│       └── main.go          # Database seeder
├── internal/
│   ├── auth/                # Authentication module
│   ├── challenges/          # Challenges module
│   ├── submissions/         # Submissions module
│   ├── middleware/          # HTTP middlewares
│   ├── migrations/          # SQL migrations (Goose)
│   ├── response/            # Response helpers
│   └── validator/           # Input validation
├── challenges/              # Challenge files (markdown, tests)
├── .env.example
├── docker-compose.yml
├── go.mod
└── Makefile
```

## API Endpoints

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/github/login` | Redirect to GitHub OAuth |
| GET | `/api/auth/github/callback` | GitHub OAuth callback |
| GET | `/api/auth/me` | Get current user (requires auth) |

### Challenges
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/challenges` | List all challenges |
| GET | `/api/challenges/:id` | Get challenge by ID |

### Submissions
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/submissions` | Create submission (requires auth) |
| GET | `/api/submissions` | Get user submissions (requires auth) |
| GET | `/api/submissions/:challenge_id` | Get submission for challenge (requires auth) |
