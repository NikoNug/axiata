package controllers

import (
	"axiata/db"
	"axiata/models"
	"encoding/json"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()
	post := new(models.Post)
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(`INSERT INTO posts (title, content, status, publish_date) VALUES ($1, $2, $3, $4) RETURNING id`, post.Title, post.Content, post.Status, post.PublishDate).Scan(&post.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, tag := range post.Tags {
		var tagID int
		err := db.QueryRow(`INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label=EXCLUDED.label RETURNING id`, tag.Label).Scan(&tagID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, post.ID, tagID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	db := db.ConnectDB()

	rows, err := db.Query(`SELECT id, title, content, status, publish_date FROM posts`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request, id int) {
	db := db.ConnectDB()
	defer db.Close()

	var post models.Post
	err := db.QueryRow(`SELECT id, title, content, status, publish_date FROM posts WHERE id = $1`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT t.id, t.label FROM tags t INNER JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id=$1`, post.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Label); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post.Tags = append(post.Tags, tag)
	}

	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, id int) {
	db := db.ConnectDB()

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`UPDATE posts SET title=$1, content=$2, status=$3, publish_date=$4 WHERE id=$5`, post.Title, post.Content, post.Status, post.PublishDate, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`DELETE FROM post_tags WHERE post_id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, tag := range post.Tags {
		var tagID int
		err := db.QueryRow(`INSERT INTO tags (label) VALUES ($1) ON CONFLICT (label) DO UPDATE SET label=EXCLUDED.label RETURNING id`, tag.Label).Scan(&tagID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(`INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`, id, tagID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request, id int) {
	db := db.ConnectDB()

	_, err := db.Exec(`DELETE FROM post_tags WHERE post_id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
