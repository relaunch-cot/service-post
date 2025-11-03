package mysql

import (
	"context"
	"fmt"
	"time"

	libModels "github.com/relaunch-cot/lib-relaunch-cot/models"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mysqlResource struct {
	mysqlClient *mysql.Client
}

type IMySqlPost interface {
	CreatePost(ctx *context.Context, userId, postId, title, content, postType, urlImagePost string) error
	GetPost(ctx *context.Context, postId string) (*libModels.Post, error)
	GetAllPosts(ctx *context.Context) ([]*libModels.Post, error)
	UpdatePost(ctx *context.Context, postId, userId, title, content, urlImagePost string) error
}

func (m *mysqlResource) CreatePost(ctx *context.Context, userId, postId, title, content, postType, urlImagePost string) error {
	currentTime := time.Now()

	if urlImagePost == "" {
		urlImagePost = "NULL"
	} else {
		urlImagePost = fmt.Sprintf("'%s'", urlImagePost)
	}
	baseQuery := fmt.Sprintf(`INSERT INTO posts (authorId, postId, title, content, type, urlImagePost, createdAt) VALUES (?, ?, ?, ?, ?, %s, ?)`, urlImagePost)
	_, err := mysql.DB.ExecContext(*ctx, baseQuery, userId, postId, title, content, postType, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (m *mysqlResource) GetPost(ctx *context.Context, postId string) (*libModels.Post, error) {
	var post libModels.Post

	query := `
SELECT 
    p.postId,
    p.authorId,
    u.name,
    p.title,
    p.content,
    p.type,
   	IFNULL(p.urlImagePost, "") AS urlImagePost,
	p.createdAt, 
	IFNULL(p.updatedAt, "") AS updatedAt
FROM posts p 
	JOIN users u ON p.authorId = u.userId
WHERE postId = ?`
	rows, err := mysql.DB.QueryContext(*ctx, query, postId)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	err = rows.Scan(
		&post.PostId,
		&post.AuthorId,
		&post.AuthorName,
		&post.Title,
		&post.Content,
		&post.Type,
		&post.UrlImagePost,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return &post, nil
}

func (m *mysqlResource) GetAllPosts(ctx *context.Context) ([]*libModels.Post, error) {
	query := `
SELECT 
	p.postId,
	p.authorId,
	u.name,
	p.title,
	p.content,
	p.type,
	IFNULL(p.urlImagePost, "") AS urlImagePost,
	p.createdAt, 
	IFNULL(p.updatedAt, "") AS updatedAt
FROM posts p 
	JOIN users u ON p.authorId = u.userId
ORDER BY p.createdAt DESC`

	rows, err := mysql.DB.QueryContext(*ctx, query)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()

	posts := make([]*libModels.Post, 0)
	for rows.Next() {
		p := &libModels.Post{}
		err = rows.Scan(
			&p.PostId,
			&p.AuthorId,
			&p.AuthorName,
			&p.Title,
			&p.Content,
			&p.Type,
			&p.UrlImagePost,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (m *mysqlResource) UpdatePost(ctx *context.Context, postId, userId, title, content, urlImagePost string) error {
	currentTime := time.Now()

	queryValidate := `
SELECT 
    p.authorId,
    p.title,
	p.content,
	IFNULL(p.urlImagePost, "") AS urlImagePost
FROM posts p 
WHERE postId = ?`

	var p libModels.Post
	rows, err := mysql.DB.QueryContext(*ctx, queryValidate, postId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return status.Error(codes.NotFound, "post not found")
	}

	err = rows.Scan(&p.AuthorId, &p.Title, &p.Content, &p.UrlImagePost)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql rows. Details: "+err.Error())
	}

	if p.AuthorId != userId {
		return status.Error(codes.PermissionDenied, "user not authorized to update this post")
	}

	var setParts []string
	if content != "" && p.Content != content {
		setParts = append(setParts, fmt.Sprintf("content = '%s'", content))
	}
	if urlImagePost != "" && p.UrlImagePost != urlImagePost {
		setParts = append(setParts, fmt.Sprintf("urlImagePost = '%s'", urlImagePost))
	}
	if title != "" && p.Title != title {
		setParts = append(setParts, fmt.Sprintf("title = '%s'", title))
	}

	if len(setParts) == 0 {
		return status.Error(codes.NotFound, "no fields to update")
	}

	setClause := setParts[0]
	for i := 1; i < len(setParts); i++ {
		setClause += ", " + setParts[i]
	}

	updateQuery := fmt.Sprintf(`UPDATE posts SET updatedAt = ?, %s WHERE postId = '%s'`, setClause, postId)

	_, err = mysql.DB.ExecContext(*ctx, updateQuery, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func NewMysqlRepository(mysqlClient *mysql.Client) IMySqlPost {
	return &mysqlResource{
		mysqlClient: mysqlClient,
	}
}
