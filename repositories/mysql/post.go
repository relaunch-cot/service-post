package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
)

type mysqlResource struct {
	mysqlClient *mysql.Client
}

type IMySqlPost interface {
	CreatePost(ctx *context.Context, userId, postId, title, content, postType, urlImagePost string) error
}

func (m *mysqlResource) CreatePost(ctx *context.Context, userId, postId, title, content, postType, urlImagePost string) error {
	currentTime := time.Now()

	if urlImagePost == "" {
		urlImagePost = "NULL"
	}
	baseQuery := fmt.Sprintf(`INSERT INTO posts (authorId, postId, title, content, type, urlImagePost, createdAt) VALUES (?, ?, ?, ?, ?, %s, ?)`, urlImagePost)
	_, err := mysql.DB.ExecContext(*ctx, baseQuery, userId, postId, title, content, postType, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func NewMysqlRepository(mysqlClient *mysql.Client) IMySqlPost {
	return &mysqlResource{
		mysqlClient: mysqlClient,
	}
}
