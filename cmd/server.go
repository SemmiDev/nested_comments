package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nested_comments "nested-comments"
	"net/http"
)

type Server struct {
	Router      *chi.Mux
	postRepo    *nested_comments.PostRepository
	commentRepo *nested_comments.CommentRepository
}

func CreateNewServer(
	postRepo *nested_comments.PostRepository,
	commentRepo *nested_comments.CommentRepository,
) *Server {
	s := &Server{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}

	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Use(middleware.Heartbeat("/ping"))

	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	s.Router.Route("/api/posts", func(postRouter chi.Router) {
		postRouter.Post("/", s.createNewPost)
		postRouter.Get("/", s.getAllPosts)
		postRouter.Get("/{id}", s.getPostByID)
		postRouter.Put("/{id}", s.updatePost)
		postRouter.Delete("/{id}", s.deletePost)

		postRouter.Post("/{id}/comments", s.createNewComment)
		postRouter.Post("/{id}/comments/{parentCommentId}", s.createNewChildComment)
		postRouter.Get("/{id}/comments", s.getParentComments)
		postRouter.Get("/{id}/comments/{parentCommentId}", s.getRepliesForParentComment)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
