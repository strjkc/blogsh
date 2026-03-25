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

func dbPostToPost(dbPost queries.Post) (handlers.Post, error) {
	createdAt, err := time.Parse(time.RFC3339, dbPost.Createdat)
	if err != nil {
		return handlers.Post{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, dbPost.Updatedat)
	if err != nil {
		return handlers.Post{}, err
	}
	return handlers.Post{
		ID:        int(dbPost.ID),
		Content:   dbPost.Content,
		Title:     dbPost.Title,
		Tags:      strings.Split(dbPost.Tags, ","),
		Category:  dbPost.Category,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (ctrl *DBController) GetPostsFromDb(queryVal string) ([]handlers.Post, error) {
	var posts []handlers.Post
	var postsDB []queries.Post
	var err error
	if len(queryVal) > 0 {
		postsDB, err = ctrl.db.GetPostsFilter(context.Background(), queries.GetPostsFilterParams{queryVal, queryVal, queryVal})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		postsDB, err = ctrl.db.GetPosts(context.Background())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	for _, dbPost := range postsDB {
		post, err := dbPostToPost(dbPost)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ctrl *DBController) GetPostFromDb(id int) (handlers.Post, error) {
	postDB, err := ctrl.db.GetPost(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}

	updatedPost, err := dbPostToPost(postDB)
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}
	return updatedPost, nil
}

func (ctrl *DBController) DeletePostFromDb(id int) error {
	_, err := ctrl.db.Delete(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (ctrl *DBController) UpdatePostInDb(id int, post handlers.Post) (handlers.Post, error) {
	origDbPost, err := ctrl.db.GetPost(context.Background(), int64(id))
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}

	updDbPost, err := ctrl.db.UpdatePost(context.Background(),
		queries.UpdatePostParams{
			ID:        origDbPost.ID,
			Title:     post.Title,
			Content:   post.Content,
			Category:  post.Category,
			Tags:      strings.Join(post.Tags, ","),
			Updatedat: time.Now().Format(time.RFC3339),
		})

	updatedPost, err := dbPostToPost(updDbPost)
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}
	return updatedPost, nil
}

func (ctrl *DBController) InsertPostIntoDb(b handlers.Post) (handlers.Post, error) {
	tagsStr := strings.Join(b.Tags, ",")
	creatTime := time.Now().Format(time.RFC3339)
	dbPost, err := ctrl.db.InsertPost(context.Background(), queries.InsertPostParams{
		Content: b.Content, Title: b.Title, Category: b.Category, Tags: tagsStr, Createdat: creatTime, Updatedat: creatTime,
	})
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}

	createdPost, err := dbPostToPost(dbPost)
	if err != nil {
		fmt.Println(err)
		return handlers.Post{}, err
	}
	return createdPost, nil
}
