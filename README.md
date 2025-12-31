# APSchool

<p align="center">
  <img src="https://img.shields.io/badge/Angular-21.0.6-dd0031" alt="Angular">
  <img src="https://img.shields.io/badge/Go-1.25.5-00add8" alt="Go">
</p>

A web platform for ESPOL students to practice Python through interactive coding challenges. Code runs directly in the browser using Pyodide (Python compiled to WebAssembly).

**Demo:** https://d3ow6eel2jvdfj.cloudfront.net

## Screenshots

| Home (logged out) | Home (logged in) |
|-------------------|------------------|
| ![Home Unauth](images/apschool_home_unauth.png) | ![Home Auth](images/apschool_home_auth.png) |

| Challenges | Challenge solved |
|------------|------------------|
| ![Challenges](images/apschool_challenges.png) | ![Test Passed](images/apschool_test_passed.png) |

| GitHub OAuth |
|--------------|
| ![GitHub OAuth](images/apschool_github_oauth.png) |

## Architecture

```
+--------------------------------------------------+
|                  BROWSER                         |
|  +------------+  +----------+  +---------------+ |
|  |  Angular   |->|  Pyodide |<-| Monaco Editor | |
|  +------------+  +----------+  +---------------+ |
|        |              |                          |
|        |         Executes Python                 |
|        |         Validates tests                 |
|        |              |                          |
|        v              v                          |
|   If tests pass -> POST /api/submissions         |
+--------------------------------------------------+
                       |
                       v
+--------------------------------------------------+
|                  AWS                             |
|  +-------------+  +-----+  +------------------+  |
|  | CloudFront  |->| EC2 |->| RDS (PostgreSQL) |  |
|  | + S3        |  | Go  |  |                  |  |
|  +-------------+  +-----+  +------------------+  |
+--------------------------------------------------+
```

## Stack

### Frontend
- Angular 21
- Pyodide

### Backend
- Go 1.25
- Chi Router
- PostgreSQL
- GitHub OAuth + JWT

### Infrastructure (AWS)
- CloudFront + S3 (frontend)
- EC2 + Docker (backend)
- RDS PostgreSQL (database)

## Features

- Python execution in the browser (no server required)
- GitHub OAuth login
- 12 challenges organized in 4 units
- Code editor with syntax highlighting
- Hint system for each challenge

## Why Pyodide?

| Aspect | Benefit |
|--------|---------|
| Cost | $0 - uses the user's CPU |
| Security | Malicious code only affects the user who wrote it |
| Latency | Instant (no round-trip to server) |
| Scalability | 1 or 100,000 users = same cost |

## Local Development

See detailed instructions at:
- [Frontend (Angular)](web/README.md)
- [Backend (Go)](server/README.md)

## Author

Johann Ramirez - johrespi@espol.edu.ec
