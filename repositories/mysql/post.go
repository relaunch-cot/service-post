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
	GetAllPostsFromUser(ctx *context.Context, userId string) ([]*libModels.Post, error)
	UpdatePost(ctx *context.Context, postId, userId, title, content, urlImagePost string) error
	DeletePost(ctx *context.Context, postId, userId string) error
	GetLikesFromPost(ctx *context.Context, postId, userId string) (*libModels.PostLikes, error)
	UpdateLikesFromPost(ctx *context.Context, postId, userId string) error
	GetAllCommentsFromPost(ctx *context.Context, postId string) (*libModels.PostComments, error)
	AddCommentToPost(ctx *context.Context, postId, commentId, userId, content string) error
	RemoveCommentFromPost(ctx *context.Context, postId, commentId, userId string) error
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
WHERE p.postId = ?`
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

func (m *mysqlResource) GetAllPostsFromUser(ctx *context.Context, userId string) ([]*libModels.Post, error) {

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
WHERE p.authorId = ?
ORDER BY p.createdAt DESC`
	rows, err := mysql.DB.QueryContext(*ctx, query, userId)
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
WHERE p.postId = ?`

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
		return status.Error(codes.PermissionDenied, "user is not authorized to perform this action")
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

func (m *mysqlResource) DeletePost(ctx *context.Context, postId, userId string) error {
	queryValidate := `
SELECT 
	p.authorId
FROM posts p 
WHERE p.postId = ?`

	var authorId string
	rows, err := mysql.DB.QueryContext(*ctx, queryValidate, postId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return status.Error(codes.NotFound, "post not found")
	}

	err = rows.Scan(&authorId)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql rows. Details: "+err.Error())
	}

	if authorId != userId {
		return status.Error(codes.PermissionDenied, "user is not authorized to perform this action")
	}

	deleteQuery := `DELETE FROM posts WHERE postId = ?`
	_, err = mysql.DB.ExecContext(*ctx, deleteQuery, postId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (m *mysqlResource) GetLikesFromPost(ctx *context.Context, postId, userId string) (*libModels.PostLikes, error) {
	postLikes := new(libModels.PostLikes)

	query := `
SELECT 
    l.userId,
    l.userName,
	l.likedAt
FROM likes l
WHERE l.postId = ?`
	rows, err := mysql.DB.QueryContext(*ctx, query, postId)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()

	var likes []libModels.Like
	for rows.Next() {
		var like libModels.Like
		err = rows.Scan(
			&like.UserId,
			&like.UserName,
			&like.LikedAt,
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
		}

		likes = append(likes, like)
	}

	postLikes.Likes = likes
	postLikes.LikesCount = int64(len(likes))

	return postLikes, nil
}

func (m *mysqlResource) UpdateLikesFromPost(ctx *context.Context, postId, userId string) error {
	currentTime := time.Now()
	var userName string

	queryValidateUser := `
SELECT 
	u.name
FROM users u 
WHERE u.userId = ?`
	rowUser, err := mysql.DB.QueryContext(*ctx, queryValidateUser, userId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rowUser.Close()
	if !rowUser.Next() {
		return status.Error(codes.NotFound, "user not found")
	}

	err = rowUser.Scan(&userName)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}

	postLikes, err := m.GetLikesFromPost(ctx, postId, userId)
	if err != nil {
		return err
	}

	isUserAlreadyLiked := false
	for _, like := range postLikes.Likes {
		if like.UserId == userId {
			isUserAlreadyLiked = true
			break
		}
	}

	if !isUserAlreadyLiked {
		queryInsert := `INSERT INTO likes (userId, postId, userName, likedAt) VALUES (?, ?, ?, ?)`
		_, err = mysql.DB.ExecContext(*ctx, queryInsert, userId, postId, userName, currentTime.Format("2006-01-02 15:04:05"))
		if err != nil {
			return status.Error(codes.Internal, "error with database. Details: "+err.Error())
		}
	} else {
		queryDelete := `DELETE FROM likes WHERE userId = ? AND postId = ?`
		_, err = mysql.DB.ExecContext(*ctx, queryDelete, userId, postId)
		if err != nil {
			return status.Error(codes.Internal, "error with database. Details: "+err.Error())
		}
	}

	return nil
}

func (m *mysqlResource) GetAllCommentsFromPost(ctx *context.Context, postId string) (*libModels.PostComments, error) {
	comments := make([]libModels.Comment, 0)

	query := `
SELECT 
	c.commentId,
	c.userId,
	c.userName,
	c.content,
	c.createdAt,
	IFNULL(c.updatedAt, "") AS updatedAt
FROM comments c 
WHERE c.postId = ?
ORDER BY c.createdAt DESC`

	rows, err := mysql.DB.QueryContext(*ctx, query, postId)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var comment libModels.Comment
		err = rows.Scan(
			&comment.CommentId,
			&comment.UserId,
			&comment.UserName,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
		}

		comments = append(comments, comment)
	}

	result := &libModels.PostComments{
		Comments:      comments,
		CommentsCount: int64(len(comments)),
	}

	return result, nil
}

func (m *mysqlResource) AddCommentToPost(ctx *context.Context, postId, commentId, userId, content string) error {
	currentTime := time.Now()
	var userName string

	queryValidateUser := `
SELECT 
	u.name
FROM users u 
WHERE u.userId = ?`

	rowUser, err := mysql.DB.QueryContext(*ctx, queryValidateUser, userId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rowUser.Close()
	if !rowUser.Next() {
		return status.Error(codes.NotFound, "user not found")
	}

	err = rowUser.Scan(&userName)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}

	baseQuery := `INSERT INTO comments (commentId, postId, userId, userName, content, createdAt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = mysql.DB.ExecContext(*ctx, baseQuery, commentId, postId, userId, userName, content, currentTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (m *mysqlResource) RemoveCommentFromPost(ctx *context.Context, postId, commentId, userId string) error {
	queryValidate := `
SELECT 
	c.userId
FROM comments c 
WHERE c.commentId = ? AND c.postId = ?`

	var commentUserId string
	rows, err := mysql.DB.QueryContext(*ctx, queryValidate, commentId, postId)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return status.Error(codes.NotFound, "comment not found")
	}

	err = rows.Scan(&commentUserId)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql rows. Details: "+err.Error())
	}

	if commentUserId != userId {
		return status.Error(codes.PermissionDenied, "user is not authorized to perform this action")
	}

	deleteQuery := `DELETE FROM comments WHERE commentId = ? AND postId = ?`
	_, err = mysql.DB.ExecContext(*ctx, deleteQuery, commentId, postId)
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
