package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jbechler2/grant-tool/backend/config"
	"github.com/jbechler2/grant-tool/backend/internal/db"
)

func main() {
	cfg := config.Load()
	fmt.Printf("grant-tool API starting, connecting to DB at %s:%s\n", cfg.DBHost, cfg.DBPort)

	database := db.Connect(cfg.DBURL)
	defer database.Close()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
