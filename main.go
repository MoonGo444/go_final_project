package main

import (
	"log"

	"github.com/moongo444/go_final_project/internal/db"
	"github.com/moongo444/go_final_project/internal/server"
)

func main() {
	if err := db.Init("TODO_DBFILE"); err != nil {
		log.Fatalf("db: %v", err)
	}

	if err := server.Run("TODO_PORT"); err != nil {
		log.Fatalf("server: %v", err)
	}
}
