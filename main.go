package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func initFrontEnd(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fs := http.FileServer(http.Dir("./frontend"))
	fs.ServeHTTP(w, r)
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	manager := NewManager(ctx)

	http.HandleFunc("/", initFrontEnd)
	http.HandleFunc("/ws", manager.ServeWS)
	http.HandleFunc("/login", manager.loginHandler)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Server started on port: 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
