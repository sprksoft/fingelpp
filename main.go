package main

import (
	"fingelpp/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/lesson/{lesson_id}", handler.Lesson)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./www/static/favicon.ico")
	})

	http.Handle("/static/", http.FileServer(http.Dir("./www")))

	port := "2025"
	fmt.Println("Starting FinGel++ HTTP server on port: " + port)
	if err := http.ListenAndServe("localhost:"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
