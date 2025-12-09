package main

import (
	"os"

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
