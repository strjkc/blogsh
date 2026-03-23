package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Unable to load environment variables")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("Unable to parse port number")
	}

	dbPath := os.Getenv("DBPATH")
	if err != nil {
		panic("Unable to fetch db path")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic("Unable to open DB file under the given path")
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	server.ListenAndServe()
	fmt.Printf("server listening on port %d\n", port)
}
