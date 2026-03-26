package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	GetPostsFromDb(queryVal string) ([]Post, error)
	GetPostFromDb(id int) (Post, error)
	DeletePostFromDb(id int) error
	UpdatePostInDb(id int, post Post) (Post, error)
	InsertPostIntoDb(b Post) (Post, error)
}

type Handler struct {
	dbcontroller Storage
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewHandler(controller Storage) Handler {
	return Handler{
		dbcontroller: controller,
	}
}

func validatePost(b Post) bool {
	if len(b.Title) < 1 || len(b.Content) < 1 || len(b.Category) < 1 || len(b.Tags) < 1 {
		return false
	}
	return true
}

func respond(w http.ResponseWriter, statusCode int, content []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(content)
}

func respondWithError(w http.ResponseWriter, statusCode int, content string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(content))
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	queryVal := r.URL.Query().Get("term")

	posts, err := h.dbcontroller.GetPostsFromDb(strings.ToLower(queryVal))
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	postsJSON, err := json.Marshal(posts)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, 200, postsJSON)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		respondWithError(w, 400, "Invalid Post")
		fmt.Println(err)
		return
	}

	post, err := h.dbcontroller.GetPostFromDb(postID)
	if err != nil {
		respondWithError(w, 404, "Post Not Found")
		fmt.Println(err)
		return
	}

	postJSON, err := json.Marshal(post)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}
	respond(w, 200, postJSON)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		respondWithError(w, 400, "Invalid Request")
		fmt.Println(err)
		return
	}

	err = h.dbcontroller.DeletePostFromDb(postID)
	if err != nil {
		respondWithError(w, 404, "Post Not Found")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusNoContent, []byte(""))
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var b Post
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		respondWithError(w, 400, "Invalid request")
		fmt.Println(err)
		return
	}
	ok := validatePost(b)
	if !ok {
		respondWithError(w, 400, "Post Post Invalid")
		return
	}

	postID, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		respondWithError(w, 400, "Invalid Request")
		fmt.Println(err)
		return
	}

	updatedPost, err := h.dbcontroller.UpdatePostInDb(postID, b)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	postJSON, err := json.Marshal(updatedPost)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusOK, postJSON)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var b Post
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		respondWithError(w, 400, "Invalid request")
		fmt.Println(err)
		return
	}

	ok := validatePost(b)
	if !ok {
		respondWithError(w, 400, "Post Post Invalid")
		return
	}

	post, err := h.dbcontroller.InsertPostIntoDb(b)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	postJSON, err := json.Marshal(post)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusCreated, postJSON)
}
