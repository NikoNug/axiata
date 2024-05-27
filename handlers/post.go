package handlers

import (
	"axiata/controllers"
	"net/http"
	"strconv"
	"strings"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		controllers.CreatePost(w, r)
	case http.MethodGet:
		controllers.GetPosts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/post/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Indalid Post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		controllers.GetPost(w, r, id)
	case http.MethodPut:
		controllers.UpdatePost(w, r, id)
	case http.MethodDelete:
		controllers.DeletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
