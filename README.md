# AP School

AP School is a web application made for students for the Programming Foundations course in ESPOL. It helps students to practice their python abilities through coding. AP School contains similar content described in the course syllabus.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Execute migrations
```bash
source .env && goose -dir internal/migrations postgres "$DATABASE_URL" up
```

## Author
Johann Ramírez - johrespi@espol.edu.ec
