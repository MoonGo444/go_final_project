package db

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date DATE NOT NULL DEFAULT "",
	title VARCHAR NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) (*sql.DB, error) {
	_ = godotenv.Load()
	dbPath := os.Getenv(dbFile)

	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	_, err := os.Stat(dbPath)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
