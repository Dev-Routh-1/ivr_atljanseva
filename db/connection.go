package db

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"github.com/joho/godotenv"


	_ "github.com/lib/pq"
)

var (
	DB   *sql.DB
	once sync.Once
)

func Connect() *sql.DB {
	_ = godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")

	once.Do(func() {
		db, err := sql.Open("postgres", dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		DB = db
		log.Println("database connected")
	})

	return DB
}