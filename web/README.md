# AP School - Frontend

Angular web application for the AP School platform.

## Requirements

- Bun 1.3+

## Setup

### 1. Install Dependencies

```bash
bun install
```

### 2. Start Development Server

```bash
bun start
```

The app will be available at `http://localhost:4200`.

## Scripts

| Command | Description |
|---------|-------------|
| `bun start` | Start development server |
| `bun run build` | Build for production |
| `bun run test` | Run tests |
| `bun run watch` | Build and watch for changes |

## Pyodide - Python in the Browser

This app uses [Pyodide](https://pyodide.org/) to execute Python code directly in the browser via WebAssembly. This means:

- **Zero server cost** for code execution - runs on user's CPU
- **Instant feedback** - no round-trip to a server
- **Secure** - malicious code only affects the user who wrote it
- **Scalable** - 1 or 100,000 users = same infrastructure cost

The `PyodideService` handles loading Pyodide from CDN and executing user code against challenge tests.

## Project Structure

```
web/
├── src/
│   ├── app/
│   │   ├── core/
│   │   │   ├── interceptors/    # HTTP interceptors
│   │   │   ├── models/          # TypeScript interfaces
│   │   │   └── services/        # API, Auth, Pyodide services
│   │   ├── features/
│   │   │   ├── auth/            # Auth callback
│   │   │   ├── challenge/       # Challenge page
│   │   │   ├── home/            # Home page
│   │   │   └── unit/            # Unit page
│   │   ├── shared/
│   │   │   ├── components/      # Reusable components
│   │   │   └── layouts/         # Layout components
│   │   └── themes/              # CSS theme tokens
│   ├── assets/                  # Static assets
│   └── environments/            # Environment config
├── angular.json
├── package.json
└── tsconfig.json
```

## Core Services

| Service | Description |
|---------|-------------|
| `ApiService` | HTTP client for backend API |
| `AuthService` | Authentication state management |
| `PyodideService` | Python execution in browser via WebAssembly |
| `ChallengeService` | Challenge data fetching |
| `SubmissionService` | Submission management |
