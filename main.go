package main

import (
	"axiata/controllers"
	"axiata/db"
	"axiata/routes"
	"log"
	"net/http"
)

func main() {
	db.ConnectDB()

	http.HandleFunc("/", controllers.CreatePost)

	router := routes.SetupRouter()

	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
