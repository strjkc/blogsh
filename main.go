package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/strjkc/blogsh/internal/database"
	"github.com/strjkc/blogsh/internal/handlers"
	"github.com/strjkc/blogsh/internal/queries"
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
	dbQueries := queries.New(db)
	controller := database.NewDBController(dbQueries)
	handler := handlers.NewHandler(&controller)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", handler.CreatePost)
	mux.HandleFunc("PUT /posts/{postID}", handler.UpdatePost)
	mux.HandleFunc("GET /posts/{postID}", handler.GetPost)
	mux.HandleFunc("DELETE /posts/{postID}", handler.DeletePost)
	mux.HandleFunc("GET /posts/", handler.GetPosts)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	fmt.Printf("server listening on port %d\n", port)
	server.ListenAndServe()
}
