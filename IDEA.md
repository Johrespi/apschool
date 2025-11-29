# Python Practice Platform

Plataforma web para practicar Python con challenges interactivos. El código se ejecuta en el navegador del usuario usando Pyodide (Python compilado a WebAssembly).

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
│         │         Ejecuta código Python                         │
│         │         Valida tests                                  │
│         │                  │                                    │
│         │                  ▼                                    │
│         │         ¿Pasó los tests?                              │
│         │            │         │                                │
│         │           SÍ        NO                                │
│         │            │         │                                │
│         │            ▼         ▼                                │
│         │      Enviar al   Mostrar error                        │
│         │       backend    (no hay request)                     │
│         │            │                                          │
└─────────┼────────────┼──────────────────────────────────────────┘
          │            │
          │            ▼
┌─────────┼────────────────────────────────────────────────────────┐
│         │       BACKEND (Go en Railway)                          │
│         │                                                        │
│         ▼                                                        │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  API REST                                                │    │
│  │  - POST /api/auth/github     (GitHub OAuth)              │    │
│  │  - GET  /api/challenges      (listar challenges)         │    │
│  │  - POST /api/submissions     (guardar solución)          │    │
│  │  - GET  /api/leaderboard     (ranking)                   │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  PostgreSQL (Supabase)                                   │    │
│  │  - users                                                 │    │
│  │  - challenges                                            │    │
│  │  - submissions                                           │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## ¿Por qué Pyodide?

Pyodide es el intérprete de Python (CPython) compilado a WebAssembly. Esto permite ejecutar Python directamente en el navegador.

| Aspecto | Beneficio |
|---------|-----------|
| Costo | $0 - usa la CPU del usuario |
| Seguridad | Código malicioso solo afecta al usuario que lo escribe |
| Latencia | Instantánea (no hay round-trip al servidor) |
| Escalabilidad | 1 o 100,000 usuarios = mismo costo |
| Simplicidad | No necesitas Docker, VMs, ni sandboxing |

### Limitaciones de Pyodide

- Primera carga: ~10-15MB (se cachea después)
- No disponible: `subprocess`, `threading` real, file system, network
- Algunas librerías no funcionan (pero numpy, pandas sí)

Para challenges de práctica (algoritmos, estructuras de datos, strings), estas limitaciones no importan.

---

## Librerías Disponibles en Pyodide

Pyodide soporta 200+ librerías. Las más relevantes:

| Categoría | Librerías |
|-----------|-----------|
| **Data Science** | pandas, numpy, scipy |
| **Machine Learning** | scikit-learn, xgboost, lightgbm |
| **Visualización** | matplotlib, bokeh, altair |
| **Matemáticas** | sympy, mpmath, statsmodels |
| **Utilidades** | regex, pyyaml, requests |

### Cómo funcionan las librerías por challenge

Cada challenge especifica qué librerías necesita en el campo `packages`:

```
packages: []           → Solo Python puro (instantáneo)
packages: ["pandas"]   → Carga pandas antes de ejecutar (2-3s primera vez)
packages: ["numpy", "scipy"] → Carga ambas librerías
```

Las librerías se cachean en el browser, así que solo la primera carga es lenta.

---

## Stack Tecnológico

### Frontend (Netlify)
```yaml
Framework: Angular 17+ (Standalone Components)
Editor: Monaco Editor
Ejecución Python: Pyodide
UI: Angular Material o Tailwind
Hosting: Railway/ Render
```

### Backend (Railway)
```yaml
Lenguaje: Go
Framework: Chi
Database: PostgreSQL (Supabase)
Auth: GitHub OAuth
Hosting: Railway / Render (gratis tier)
```

### Servicios Externos
```yaml
Database: Supabase (PostgreSQL gratis)
Auth: GitHub OAuth (gratis)
CDN Pyodide: jsDelivr (gratis)
```

---

## Flujo de Usuario

```
1. Usuario visita la app (Netlify sirve Angular)
2. Login con GitHub OAuth
3. Ve lista de challenges
4. Selecciona un challenge
5. Escribe código en Monaco Editor
6. Click "Ejecutar"
   └── Pyodide ejecuta Python EN EL BROWSER
   └── Valida contra los tests del challenge
7. Si pasa:
   └── POST /api/submissions al backend
   └── Backend guarda en PostgreSQL
   └── Actualiza leaderboard
8. Si falla:
   └── Muestra error (no hay request al backend)
```

---

## Base de Datos

### Tabla: users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    github_id INTEGER UNIQUE NOT NULL,
    username VARCHAR(50) NOT NULL,
    avatar_url VARCHAR(500),
    score INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Tabla: challenges
```sql
CREATE TABLE challenges (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20), -- 'easy', 'medium', 'hard'
    category VARCHAR(50), -- 'basics', 'data-science', 'algorithms'
    packages TEXT[], -- ['pandas', 'numpy'] o []
    template TEXT NOT NULL, -- código inicial
    test_code TEXT NOT NULL, -- tests que valida Pyodide
    hints TEXT[],
    points INTEGER DEFAULT 10,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Tabla: submissions
```sql
CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    challenge_id INTEGER REFERENCES challenges(id),
    code TEXT NOT NULL,
    passed BOOLEAN NOT NULL,
    execution_time FLOAT, -- en ms (medido en el browser)
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, challenge_id) -- solo una solución por challenge
);
```

---

## API Endpoints

### Auth
```
GET  /api/auth/github          - Redirige a GitHub OAuth
GET  /api/auth/github/callback - Callback de GitHub
GET  /api/auth/me              - Usuario actual
POST /api/auth/logout          - Cerrar sesión
```

### Challenges
```
GET  /api/challenges           - Listar todos los challenges
GET  /api/challenges/:slug     - Ver un challenge
```

### Submissions
```
POST /api/submissions          - Guardar solución exitosa
GET  /api/submissions          - Mis submissions
```

### Leaderboard
```
GET  /api/leaderboard          - Top usuarios por score
```

---

## Estructura del Proyecto

```
python-practice/
├── frontend/                      # Angular (Netlify)
│   ├── src/
│   │   ├── app/
│   │   │   ├── services/
│   │   │   │   ├── pyodide.service.ts    # Ejecutar Python
│   │   │   │   ├── auth.service.ts       # GitHub OAuth
│   │   │   │   └── api.service.ts        # HTTP al backend
│   │   │   ├── components/
│   │   │   │   ├── challenge-list/
│   │   │   │   ├── challenge-detail/
│   │   │   │   ├── code-editor/
│   │   │   │   └── leaderboard/
│   │   │   └── app.routes.ts
│   │   └── index.html                    # Carga Pyodide CDN
│   └── package.json
│
├── backend/                       # Go (Railway)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── challenges.go
│   │   │   ├── submissions.go
│   │   │   └── leaderboard.go
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── challenge.go
│   │   │   └── submission.go
│   │   ├── database/
│   │   │   └── postgres.go
│   │   └── middleware/
│   │       └── auth.go
│   ├── go.mod
│   └── Dockerfile
│
└── challenges/                    # Challenges en JSON (se cargan a DB)
    ├── 001-hello-world.json
    ├── 002-two-sum.json
    └── ...
```

---

## Ejemplo: Challenges JSON

### Challenge sin librerías (básico)
```json
{
  "slug": "two-sum",
  "title": "Two Sum",
  "description": "Dada una lista de números y un target, retorna los índices de los dos números que suman el target.",
  "difficulty": "easy",
  "category": "basics",
  "packages": [],
  "points": 10,
  "template": "def two_sum(nums: list, target: int) -> list:\n    # Tu código aquí\n    pass",
  "test_code": "assert two_sum([2,7,11,15], 9) == [0,1], 'Test 1 falló'\nassert two_sum([3,2,4], 6) == [1,2], 'Test 2 falló'\nassert two_sum([3,3], 6) == [0,1], 'Test 3 falló'\nprint('ALL_TESTS_PASSED')",
  "hints": [
    "Puedes usar un diccionario para guardar los números que ya viste",
    "Para cada número, calcula el complemento (target - num)"
  ]
}
```

### Challenge con Pandas (data science)
```json
{
  "slug": "rolling-mean",
  "title": "Media Móvil",
  "description": "Calcula la media móvil de 3 días para una lista de precios.",
  "difficulty": "medium",
  "category": "data-science",
  "packages": ["pandas"],
  "points": 20,
  "template": "import pandas as pd\n\ndef rolling_mean(prices: list) -> list:\n    # Tu código aquí\n    pass",
  "test_code": "result = rolling_mean([1, 2, 3, 4, 5])\nassert len(result) == 5, 'Debe retornar 5 elementos'\nassert result[2] == 2.0, 'El tercer elemento debe ser 2.0'\nassert result[4] == 4.0, 'El quinto elemento debe ser 4.0'\nprint('ALL_TESTS_PASSED')",
  "hints": [
    "Usa pd.Series(prices).rolling(3).mean()",
    "Convierte el resultado a lista con .tolist()"
  ]
}
```

### Challenge con NumPy (algoritmos numéricos)
```json
{
  "slug": "matrix-multiply",
  "title": "Multiplicación de Matrices",
  "description": "Multiplica dos matrices usando NumPy.",
  "difficulty": "medium",
  "category": "algorithms",
  "packages": ["numpy"],
  "points": 15,
  "template": "import numpy as np\n\ndef matrix_multiply(a: list, b: list) -> list:\n    # Tu código aquí\n    pass",
  "test_code": "result = matrix_multiply([[1, 2], [3, 4]], [[5, 6], [7, 8]])\nassert result == [[19, 22], [43, 50]], 'Multiplicación incorrecta'\nprint('ALL_TESTS_PASSED')",
  "hints": [
    "Usa np.array() para convertir las listas",
    "Usa np.dot() o el operador @ para multiplicar"
  ]
}
```

---

## PyodideService (Angular)

```typescript
import { Injectable } from '@angular/core';

declare global {
  interface Window {
    loadPyodide: () => Promise<any>;
  }
}

@Injectable({ providedIn: 'root' })
export class PyodideService {
  private pyodide: any = null;
  private loading: Promise<any> | null = null;
  private loadedPackages: Set<string> = new Set();

  async init(): Promise<void> {
    if (this.pyodide) return;
    if (this.loading) {
      await this.loading;
      return;
    }
    this.loading = window.loadPyodide();
    this.pyodide = await this.loading;
  }

  async loadPackages(packages: string[]): Promise<void> {
    const toLoad = packages.filter(pkg => !this.loadedPackages.has(pkg));
    if (toLoad.length === 0) return;
    
    await this.pyodide.loadPackage(toLoad);
    toLoad.forEach(pkg => this.loadedPackages.add(pkg));
  }

  async runWithTests(
    userCode: string, 
    testCode: string,
    packages: string[] = []
  ): Promise<TestResult> {
    await this.init();
    
    // Cargar librerías si el challenge las requiere
    if (packages.length > 0) {
      await this.loadPackages(packages);
    }

    const fullCode = `
${userCode}

${testCode}
`;

    const startTime = performance.now();
    
    try {
      // Capturar stdout
      this.pyodide.runPython(`
import sys
from io import StringIO
sys.stdout = StringIO()
      `);

      this.pyodide.runPython(fullCode);
      const output = this.pyodide.runPython('sys.stdout.getvalue()');
      const executionTime = performance.now() - startTime;

      return {
        passed: output.includes('ALL_TESTS_PASSED'),
        output,
        error: null,
        executionTime
      };
    } catch (error: any) {
      return {
        passed: false,
        output: '',
        error: error.message,
        executionTime: performance.now() - startTime
      };
    }
  }
}

interface TestResult {
  passed: boolean;
  output: string;
  error: string | null;
  executionTime: number;
}
```

---

## index.html (cargar Pyodide)

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Python Practice</title>
  <script src="https://cdn.jsdelivr.net/pyodide/v0.29.0/full/pyodide.js"></script>
</head>
<body>
  <app-root></app-root>
</body>
</html>
```

---

## Costos

| Servicio | Costo |
|----------|-------|
| Netlify (frontend) | $0 |
| Railway (backend) | $0 (free tier: 500 hrs/mes) |
| Supabase (PostgreSQL) | $0 (free tier: 500MB) |
| GitHub OAuth | $0 |
| Pyodide CDN | $0 |
| **Total** | **$0/mes** |

---

## MVP - Funcionalidades

### Fase 1 (MVP)
- [ ] Login con GitHub
- [ ] Listar challenges
- [ ] Editor de código (Monaco)
- [ ] Ejecutar código con Pyodide
- [ ] Validar tests
- [ ] Guardar submissions
- [ ] Leaderboard básico

### Fase 2
- [ ] Perfil de usuario
- [ ] Filtrar challenges por dificultad
- [ ] Hints
- [ ] Historial de submissions
- [ ] Tema oscuro/claro

### Fase 3
- [ ] Más challenges
- [ ] Badges/logros
- [ ] Estadísticas de usuario
- [ ] PWA (offline)

---

## Desarrollo Local

### Frontend
```bash
cd frontend
npm install
npm start
# http://localhost:4200
```

### Backend
```bash
cd backend
go mod tidy
go run cmd/server/main.go
# http://localhost:8080
```

### Variables de Entorno

```bash
# Backend (.env)
DATABASE_URL=postgres://user:pass@host:5432/dbname
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx
JWT_SECRET=xxx
FRONTEND_URL=http://localhost:4200
```

---

## Deployment

### Frontend → Netlify
1. Conectar repo a Netlify
2. Build command: `npm run build`
3. Publish directory: `dist/frontend`

### Backend → Railway
1. Conectar repo a Railway
2. Railway detecta Go automáticamente
3. Configurar variables de entorno

### Database → Supabase
1. Crear proyecto en Supabase
2. Copiar connection string
3. Ejecutar migraciones

---

## Seguridad

| Riesgo | Solución |
|--------|----------|
| Loop infinito | Corre en el browser del usuario, no afecta al servidor |
| Código malicioso | Sandbox de WebAssembly, no puede acceder al sistema |
| Inyección SQL | Prepared statements en Go |
| XSS | Angular sanitiza por defecto |
| CSRF | Tokens JWT en header Authorization |

---

## Resumen

Esta arquitectura es simple, escalable y **completamente gratis**:

1. **Angular** sirve la UI desde Netlify
2. **Pyodide** ejecuta Python en el browser (cero costo de servidor)
3. **Go** maneja auth y persistencia en Railway
4. **Supabase** almacena usuarios, challenges y submissions

El código Python nunca toca tu servidor. Solo guardas: "usuario X completó challenge Y".
