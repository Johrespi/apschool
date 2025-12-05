# AP School - Python Practice Platform

Plataforma web para practicar Python con challenges interactivos. El codigo se ejecuta en el navegador del usuario usando Pyodide (Python compilado a WebAssembly).

---

## Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│                    BROWSER DEL USUARIO                          │
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌────────────────────┐  │
│  │   Angular   │    │   Pyodide   │    │   Monaco Editor    │  │
│  │    (UI)     │───▶│  (Python)   │◀───│   (Code Editor)    │  │
│  └─────────────┘    └─────────────┘    └────────────────────┘  │
│         │                  │                                    │
│         │                  ▼                                    │
│         │         Ejecuta codigo Python                         │
│         │         Valida tests                                  │
│         │                  │                                    │
│         │                  ▼                                    │
│         │         Paso los tests?                               │
│         │            │         │                                │
│         │           SI        NO                                │
│         │            │         │                                │
│         │            ▼         ▼                                │
│         │      Enviar al   Mostrar error                        │
│         │       backend    (no hay request)                     │
│         │            │                                          │
└─────────┼────────────┼──────────────────────────────────────────┘
          │            │
          │            ▼
┌─────────┼────────────────────────────────────────────────────────┐
│         │       BACKEND (Go)                                     │
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  API REST                                                │    │
│  │  - GET  /api/auth/github/login     (GitHub OAuth)        │    │
│  │  - GET  /api/auth/github/callback  (GitHub callback)     │    │
│  │  - GET  /api/challenges?category=  (listar challenges)   │    │
│  │  - GET  /api/challenges/:id        (obtener challenge)   │    │
│  │  - POST /api/submissions           (guardar solucion)    │    │
│  │  - GET  /api/submissions           (mis submissions)     │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  PostgreSQL                                              │    │
│  │  - users                                                 │    │
│  │  - user_auth_github                                      │    │
│  │  - user_auth_email (futuro)                              │    │
│  │  - challenges                                            │    │
│  │  - submissions                                           │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## Stack Tecnologico

### Frontend
- **Framework**: Angular 17+ (Standalone Components)
- **Editor**: Monaco Editor
- **Ejecucion Python**: Pyodide
- **UI**: Angular Material o Tailwind

### Backend
- **Lenguaje**: Go
- **Router**: Chi
- **Database**: PostgreSQL
- **Auth**: GitHub OAuth (email/password futuro)
- **Migraciones**: Goose

---

## Por que Pyodide?

Pyodide es el interprete de Python (CPython) compilado a WebAssembly. Ejecuta Python directamente en el navegador.

| Aspecto | Beneficio |
|---------|-----------|
| Costo | $0 - usa la CPU del usuario |
| Seguridad | Codigo malicioso solo afecta al usuario que lo escribe |
| Latencia | Instantanea (no hay round-trip al servidor) |
| Escalabilidad | 1 o 100,000 usuarios = mismo costo |

### Librerias en Pyodide

Pyodide detecta automaticamente los imports y carga las librerias necesarias:

```javascript
// El frontend detecta y carga librerias antes de ejecutar
await pyodide.loadPackagesFromImports(userCode);
await pyodide.runPython(userCode);
```

Librerias soportadas: numpy, pandas, scipy, scikit-learn, matplotlib, etc.

---

## Base de Datos

### Tabla: users
```sql
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    avatar_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Tabla: user_auth_github
```sql
CREATE TABLE user_auth_github (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    github_id BIGINT UNIQUE NOT NULL
);
```

### Tabla: user_auth_email (futuro)
```sql
CREATE TABLE user_auth_email (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL
);
```

### Tabla: challenges
```sql
CREATE TABLE challenges (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    template TEXT NOT NULL,
    test_code TEXT NOT NULL,
    hints TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Tabla: submissions
```sql
CREATE TABLE submissions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    passed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, challenge_id)
);
```

---

## API Endpoints

### Auth
```
GET  /api/auth/github/login     - Redirige a GitHub OAuth
GET  /api/auth/github/callback  - Callback de GitHub, retorna JWT
```

### Challenges
```
GET  /api/challenges?category=  - Listar challenges por categoria
GET  /api/challenges/:id        - Obtener un challenge completo
```

### Submissions
```
POST /api/submissions           - Guardar solucion exitosa
GET  /api/submissions           - Mis submissions
```

---

## Estructura de Challenges (Archivos)

Los challenges se organizan en carpetas por unidad:

```
challenges/
└── unit-1-intro/
    └── 001-hello-world/
        ├── README.md       # Descripcion del challenge
        ├── template.py     # Codigo inicial (incluye imports si necesita librerias)
        ├── tests.py        # Tests de validacion
        └── hints.md        # Pistas para el estudiante
```

### Ejemplo: README.md
```markdown
# Hello World

## Descripcion
Escribe un programa que imprima "Hello, World!" en la consola.

## Ejemplo
### Output esperado
Hello, World!

## Instrucciones
1. Usa la funcion `print()` para mostrar el mensaje.
```

### Ejemplo: template.py
```python
# Escribe tu codigo aqui

```

### Ejemplo: tests.py
```python
import sys
from io import StringIO

captured_output = StringIO()
sys.stdout = captured_output

# El codigo del usuario se ejecuta ANTES de este archivo
output = captured_output.getvalue().strip()

assert output == "Hello, World!", f"Se esperaba 'Hello, World!' pero se obtuvo '{output}'"
print("ALL_TESTS_PASSED")
```

### Ejemplo: hints.md
```markdown
# Hints

## Hint 1
La funcion `print()` muestra texto en la consola.

## Hint 2
Los strings en Python se escriben entre comillas: `"texto"` o `'texto'`.
```

---

## Estructura del Proyecto

```
apschool/
├── cmd/
│   └── api/
│       ├── main.go              # Entry point
│       └── routes.go            # Definicion de rutas
├── internal/
│   ├── auth/                    # Modulo de autenticacion
│   │   ├── handler.go
│   │   ├── jwt.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── challenges/              # Modulo de challenges
│   │   ├── handler.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── middleware/
│   │   └── auth.go              # JWT middleware
│   ├── migrations/              # Migraciones SQL (Goose)
│   │   ├── 001_create_users.sql
│   │   ├── 002_create_user_auth_github.sql
│   │   ├── 003_create_user_auth_email.sql
│   │   ├── 004_create_challenges.sql
│   │   └── 005_create_submissions.sql
│   ├── response/
│   │   ├── errors.go
│   │   └── json.go
│   └── validator/
│       └── validator.go
├── challenges/                  # Challenges en archivos
│   └── unit-1-intro/
│       └── 001-hello-world/
├── .air.toml                    # Hot reload config
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Flujo de Usuario

```
1. Usuario visita la app
2. Login con GitHub OAuth
3. Ve lista de challenges por categoria
4. Selecciona un challenge
5. Escribe codigo en Monaco Editor
6. Click "Ejecutar"
   └── Pyodide carga librerias si hay imports
   └── Pyodide ejecuta Python EN EL BROWSER
   └── Valida contra los tests del challenge
7. Si pasa:
   └── POST /api/submissions al backend
   └── Backend guarda en PostgreSQL
8. Si falla:
   └── Muestra error (no hay request al backend)
```

---

## MVP Checklist

### Backend
- [x] Estructura del proyecto (cmd/api, internal/)
- [x] Configuracion de Chi router
- [x] Migraciones de base de datos
- [x] Response helpers (JSON, errors)
- [x] Validator
- [x] Auth: GitHub OAuth login/callback
- [x] Auth: JWT generation y middleware
- [x] Challenges: Model, Repository, Service, Handler
- [ ] Challenges: Conectar rutas en routes.go
- [ ] Submissions: Model, Repository, Service, Handler
- [ ] Submissions: Conectar rutas

### Frontend (Angular)
- [ ] Setup proyecto Angular
- [ ] PyodideService
- [ ] AuthService (GitHub OAuth)
- [ ] Monaco Editor integration
- [ ] Challenge list component
- [ ] Challenge detail component
- [ ] Submissions history

---

## Desarrollo Local

### Requisitos
- Go 1.21+
- Docker (para PostgreSQL)

### Backend
```bash
# Iniciar base de datos
make docker-run

# Ejecutar migraciones
source .env && goose -dir internal/migrations postgres "$DATABASE_URL" up

# Iniciar servidor con hot reload
make watch

# O sin hot reload
make run
```

### Variables de Entorno (.env)
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/apschool?sslmode=disable
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
JWT_SECRET=xxx
FRONTEND_URL=http://localhost:4200
```
