package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/strjkc/blogsh/internal/handlers"
	"github.com/strjkc/blogsh/internal/queries"
)

type DBController struct {
	db *queries.Queries
}

func NewDBController(dbQueries *queries.Queries) DBController {
	return DBController{
		db: dbQueries,
	}
}

func dbBlogToBlog(dbBlog queries.Blog) (handlers.Blog, error) {
	createdAt, err := time.Parse(time.RFC3339, dbBlog.Createdat)
	if err != nil {
		return handlers.Blog{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, dbBlog.Updatedat)
	if err != nil {
		return handlers.Blog{}, err
	}
	return handlers.Blog{
		ID:        int(dbBlog.ID),
		Content:   dbBlog.Content,
		Title:     dbBlog.Title,
		Tags:      strings.Split(dbBlog.Tags, ","),
		Category:  dbBlog.Category,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (ctrl *DBController) GetBlogsFromDb(queryVal string) ([]handlers.Blog, error) {
	var blogs []handlers.Blog
	var blogsDB []queries.Blog
	var err error
	if len(queryVal) > 0 {
		blogsDB, err = ctrl.db.GetBlogsFilter(context.Background(), queries.GetBlogsFilterParams{queryVal, queryVal, queryVal})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		blogsDB, err = ctrl.db.GetBlogs(context.Background())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	for _, dbBlog := range blogsDB {
		blog, err := dbBlogToBlog(dbBlog)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (ctrl *DBController) GetBlogFromDb(id int) (handlers.Blog, error) {
	blogDB, err := ctrl.db.GetBlog(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}

	updatedBlog, err := dbBlogToBlog(blogDB)
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}
	return updatedBlog, nil
}

func (ctrl *DBController) DeleteBlogFromDb(id int) error {
	_, err := ctrl.db.Delete(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (ctrl *DBController) UpdateBlogInDb(id int, blog handlers.Blog) (handlers.Blog, error) {
	origDbBlog, err := ctrl.db.GetBlog(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}

	updDbBlog, err := ctrl.db.UpdateBlog(context.Background(),
		queries.UpdateBlogParams{
			ID:        origDbBlog.ID,
			Title:     blog.Title,
			Content:   blog.Content,
			Tags:      strings.Join(blog.Tags, ","),
			Updatedat: time.Now().Format(time.RFC3339),
		})

	updatedBlog, err := dbBlogToBlog(updDbBlog)
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}
	return updatedBlog, nil
}

func (ctrl *DBController) InsertBlogIntoDb(b handlers.Blog) (handlers.Blog, error) {
	tagsStr := strings.Join(b.Tags, ",")
	creatTime := time.Now().Format(time.RFC3339)
	dbBlog, err := ctrl.db.InsertBlog(context.Background(), queries.InsertBlogParams{
		Content: b.Content, Title: b.Title, Category: b.Category, Tags: tagsStr, Createdat: creatTime, Updatedat: creatTime,
	})
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}

	createdBlog, err := dbBlogToBlog(dbBlog)
	if err != nil {
		fmt.Println(err)
		return handlers.Blog{}, err
	}
	return createdBlog, nil
}
