package routes

import (
	"axiata/handlers"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/posts", handlers.PostsHandler)
	mux.HandleFunc("/api/post/", handlers.PostHandler)

	return mux
}
