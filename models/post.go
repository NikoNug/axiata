package models

import "time"

type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      string    `json:"status"`
	PublishDate time.Time `json:"publish_date"`
}

type Tag struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
