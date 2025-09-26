package main

import (
	"fingelpp/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.DynamicHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "content/static/favicon.ico")
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./content/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./content/js"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./content/static"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./content/media"))))

	port := "2025"
	fmt.Println("Starting FinGel++ HTTP server on port: " + port)
	if err := http.ListenAndServe("localhost:"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
