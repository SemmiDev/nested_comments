package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	ns "nested-comments"
	"net/http"
)

func (s *Server) createNewPost(w http.ResponseWriter, r *http.Request) {
	var post ns.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	param := ns.CreatePostParam{
		ID:      uuid.New(),
		Title:   post.Title,
		Content: post.Content,
	}

	err := s.postRepo.CreatePost(r.Context(), param)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(param)
}

func (s *Server) getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.postRepo.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (s *Server) getPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	comments, err := s.commentRepo.GetParentComments(post.ID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	post.Comments = comments

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (s *Server) updatePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post ns.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	param := ns.UpdatePostParam{
		Title:   post.Title,
		Content: post.Content,
	}

	if err := s.postRepo.UpdatePost(postID, param); err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deletePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := s.postRepo.DeletePost(postID); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) createNewComment(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	param := ns.CreateCommentParam{
		ID:       uuid.New(),
		ParentID: uuid.NullUUID{},
		PostID:   postID,
		Text:     comment.Text,
	}

	err = s.commentRepo.CreateComment(param)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(param)
}

func (s *Server) createNewChildComment(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	parentCommentIDStr := chi.URLParam(r, "parentCommentId")
	parentCommentID, err := uuid.Parse(parentCommentIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	param := ns.CreateCommentParam{
		ID:       uuid.New(),
		ParentID: uuid.NullUUID{UUID: parentCommentID, Valid: true},
		PostID:   postID,
		Text:     comment.Text,
	}

	err = s.commentRepo.CreateComment(param)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(param)
}

func (s *Server) getParentComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := s.commentRepo.GetParentComments(postID)
	if err != nil {
		http.Error(w, "Failed to get parent comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (s *Server) getRepliesForParentComment(w http.ResponseWriter, r *http.Request) {
	parentCommentIDStr := chi.URLParam(r, "parentCommentId")
	parentCommentID, err := uuid.Parse(parentCommentIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	replies, err := s.commentRepo.GetRepliesForParentComment(parentCommentID)
	if err != nil {
		http.Error(w, "Failed to get replies for parent comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}
