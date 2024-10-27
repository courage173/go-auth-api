package main

import (
	"log"

	"github.com/courage173/go-auth-api/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
    dburl := utils.GetEnv("DB_URL")

    if dburl == "" {
        panic("DB_URL is not set")
    }
    
    m, err := migrate.New(
		"file://migrations",
		dburl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
    m.Up() 
}