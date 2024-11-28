package main

import (
	"log"
	"net/http"
	"go-backend-technical-test/database"
	"go-backend-technical-test/routes"
)

func main() {
	database.Connect()
	r := routes.AppRoutes()

	http.Handle("/", routes.CORS(r))

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

