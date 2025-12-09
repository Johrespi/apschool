package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database = os.Getenv("APSCHOOL_DB_DATABASE")
	password = os.Getenv("APSCHOOL_DB_PASSWORD")
	username = os.Getenv("APSCHOOL_DB_USERNAME")
	host     = os.Getenv("APSCHOOL_DB_HOST")
	port     = os.Getenv("APSCHOOL_DB_PORT")
	schema   = os.Getenv("APSCHOOL_DB_SCHEMA")
)

type Challenge struct {
	Slug        string
	Category    string
	Title       string
	Description string
	Template    string
	TestCode    string
	Hints       string
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func loadChallenges(basePath string) ([]Challenge, error) {
	var challenges []Challenge

	categories, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if !category.IsDir() {
			continue
		}

		categoryPath := filepath.Join(basePath, category.Name())
		slugs, err := os.ReadDir(categoryPath)
		if err != nil {
			return nil, err
		}

		for _, slug := range slugs {
			if !slug.IsDir() {
				continue
			}

			challengePath := filepath.Join(categoryPath, slug.Name())
			challenge, err := loadChallenge(challengePath, category.Name(), slug.Name())
			if err != nil {
				log.Printf("Warning: skipping %s: %v", challengePath, err)
				continue
			}

			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func loadChallenge(path, category, slug string) (Challenge, error) {
	readme, err := readFile(filepath.Join(path, "README.md"))
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
