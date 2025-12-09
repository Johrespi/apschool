package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	if err != nil {
		return Challenge{}, fmt.Errorf("README.md :%w", err)
	}

	title, description := parseReadme(readme)

	template, err := readFile(file.Join(path, "template.py"))
	if err != nil {
		return Challenge{}, fmt.Errorf("template.py :%w", err)
	}

	testCode, err := readFile(filepath.Join(path, "tests.py"))
	if err != nil {
		return Challenge{}, fmt.Errorf("tests.py :%w", err)
	}

	hints, err := readFile(filepath.Join(path, "hints.md"))
	if err != nil {
		return Challenge{}, fmt.Errorf("hints.md :%w", err)
	}

	return Challenge{
		Slug:        slug,
		Category:    category,
		Title:       title,
		Description: description,
		Template:    template,
		TestCode:    testCode,
		Hints:       hints,
	}, nil

}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func parseReadme(content string) (title, description string) {
	lines := strings.SplitN(content, "\n", 2)

	if len(lines) == 0 {
		return "", ""
	}

	// Extract title from the first line
	title = strings.TrimPrefix(lines[0], "# ")
	// Remove any extra whitespace at the beginning and end
	title = strings.TrimSpace(title)

	if len(lines) > 1 {
		description = strings.TrimSpace(lines[1])
	}

	return title, description
}

func upsertChallenge(db *sql.DB, c Challenge) error {
	query := `
	INSERT INTO challenges(slug, category, title, description, template, test_code, hints)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (slug) DO UPDATE SET
		category = $2,
		title = $3,
		description = $4,
		template = $5,
		test_code = $6,
		hints = $7;
		updated_at = NOW()
	`

	_, err := db.Exec(query, c.Slug, c.Category, c.Title, c.Description, c.Template, c.TestCode, c.Hints)
	return err
}
