package main

import (
	"database/sql"
	"log"

	_ "github.com/bcantrell1/pro-motocross-api/docs"

	"github.com/bcantrell1/pro-motocross-api/internal/database"
	"github.com/bcantrell1/pro-motocross-api/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title Pro Motocross Rest API
// @version 1.0
// @description This is a rest API written in Go utilizing the Gin framework.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// Apply the security definition to your endpoints
// @security BearerAuth

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal((err))
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "my-super-secret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
