package postgres

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func NewDBConnect() (*sql.DB, error) {

	godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbAddr := os.Getenv("DB_IP")
	dbPort := os.Getenv("DB_PORT")

	var err error
	Db, err = sql.Open("postgres", "postgres://"+dbUser+":"+dbPass+"@"+dbAddr+":"+dbPort+"/"+dbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	Db.SetMaxOpenConns(60)
	Db.SetConnMaxLifetime(120)
	Db.SetMaxIdleConns(30)
	Db.SetConnMaxIdleTime(20)
	if Db.Ping(); err != nil {
		log.Fatal(err)
	}

	return Db, nil
}
