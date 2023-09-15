package nested_comments

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`

	Comments []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        uuid.UUID     `json:"id"`
	ParentID  uuid.NullUUID `json:"parent_id"`
	PostID    uuid.UUID     `json:"post_id"`
	Text      string        `json:"text"`
	CreatedAt string        `json:"created_at"`
}

type PostRepository struct {
	db *sql.DB
}

type CommentRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db}
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db}
}

type CreatePostParam struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

func (pr *PostRepository) CreatePost(ctx context.Context, param CreatePostParam) error {
	query := "INSERT INTO posts (id, title, content) VALUES (?, ?, ?)"
	_, err := pr.db.ExecContext(ctx, query, param.ID, param.Title, param.Content)
	return err
}

func (pr *PostRepository) GetPostByID(postID uuid.UUID) (Post, error) {
	query := "SELECT id, title, content, created_at FROM posts WHERE id = ?"
	var post Post
	err := pr.db.QueryRow(query, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	return post, err
}

type UpdatePostParam struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (pr *PostRepository) UpdatePost(postID uuid.UUID, param UpdatePostParam) error {
	query := "UPDATE posts SET title = ?, content = ? WHERE id = ?"
	_, err := pr.db.Exec(query, param.Title, param.Content, postID)
	return err
}

func (pr *PostRepository) DeletePost(postID uuid.UUID) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := pr.db.Exec(query, postID)
	return err
}

func (pr *PostRepository) GetAllPosts() ([]Post, error) {
	query := "SELECT id, title, content, created_at FROM posts"
	rows, err := pr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

type CreateCommentParam struct {
	ID       uuid.UUID     `json:"id"`
	ParentID uuid.NullUUID `json:"parent_id"`
	PostID   uuid.UUID     `json:"post_id"`
	Text     string        `json:"text"`
}

func (cr *CommentRepository) CreateComment(param CreateCommentParam) error {
	query := "INSERT INTO comments (id, parent_id, post_id, text) VALUES (?, ?, ?, ?)"
	_, err := cr.db.Exec(query, param.ID, param.ParentID, param.PostID, param.Text)
	return err
}

func (cr *CommentRepository) GetCommentByID(commentID uuid.UUID) (Comment, error) {
	query := "SELECT id, parent_id, post_id, text, created_at FROM comments WHERE id = ?"
	var comment Comment
	err := cr.db.QueryRow(query, commentID).Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func (cr *CommentRepository) GetParentComments(postID uuid.UUID) ([]Comment, error) {
	query := "SELECT id, parent_id, post_id, text, created_at FROM comments WHERE post_id = ? AND parent_id IS NULL"
	rows, err := cr.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.Text, &comment.CreatedAt)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (cr *CommentRepository) GetRepliesForParentComment(parentID uuid.UUID) ([]Comment, error) {
	query := "SELECT id, parent_id, post_id, text, created_at FROM comments WHERE parent_id = ?"
	rows, err := cr.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.ParentID, &comment.PostID, &comment.Text, &comment.CreatedAt)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
