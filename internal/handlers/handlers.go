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
	GetBlogsFromDb(queryVal string) ([]Blog, error)
	GetBlogFromDb(id int) (Blog, error)
	DeleteBlogFromDb(id int) error
	UpdateBlogInDb(id int, blog Blog) (Blog, error)
	InsertBlogIntoDb(b Blog) (Blog, error)
}

type Handler struct {
	dbcontroller Storage
}

type Blog struct {
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

func validateBlog(b Blog) bool {
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

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("<p>%s</p>", message)))
}

func (h *Handler) GetBlogs(w http.ResponseWriter, r *http.Request) {
	queryVal := r.URL.Query().Get("term")

	blogs, err := h.dbcontroller.GetBlogsFromDb(strings.ToLower(queryVal))
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	blogsJSON, err := json.Marshal(blogs)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, 200, blogsJSON)
}

func (h *Handler) GetBlog(w http.ResponseWriter, r *http.Request) {
	blogID, err := strconv.Atoi(r.PathValue("blogID"))
	if err != nil {
		respondWithError(w, 400, "Invalid blog")
		fmt.Println(err)
		return
	}

	blog, err := h.dbcontroller.GetBlogFromDb(blogID)
	if err != nil {
		respondWithError(w, 404, "Blog Not Found")
		fmt.Println(err)
		return
	}

	blogJSON, err := json.Marshal(blog)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}
	respond(w, 200, blogJSON)
}

func (h *Handler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	blogID, err := strconv.Atoi(r.PathValue("blogID"))
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	err = h.dbcontroller.DeleteBlogFromDb(blogID)
	if err != nil {
		respondWithError(w, 404, "Blog Not Found")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusNoContent, []byte(""))
}

func (h *Handler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	var b Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}
	ok := validateBlog(b)
	if !ok {
		respondWithError(w, 400, "Blog Post Invalid")
		return
	}

	blogID, err := strconv.Atoi(r.PathValue("blogID"))
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	updatedBlog, err := h.dbcontroller.UpdateBlogInDb(blogID, b)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	blogJSON, err := json.Marshal(updatedBlog)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusCreated, blogJSON)
}

func (h *Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var b Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	ok := validateBlog(b)
	if !ok {
		respondWithError(w, 400, "Blog Post Invalid")
		return
	}

	blog, err := h.dbcontroller.InsertBlogIntoDb(b)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	blogJSON, err := json.Marshal(blog)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		fmt.Println(err)
		return
	}

	respond(w, http.StatusCreated, blogJSON)
}
