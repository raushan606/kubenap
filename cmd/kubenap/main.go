package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("🚀 Starting KubeNap Controller")

	http.HandleFunc("/wake", func(w http.ResponseWriter, r *http.Request) {
		original := r.URL.Query().Get("original")
		if original == "" {
			http.Error(w, "missing original path", http.StatusBadRequest)
			return
		}

		// For now, just simulate resume logic
		log.Printf("[wake] Received request to wake app for: %s", original)
		// TODO: Implement actual resume logic here (scale up deployment)

		// Simulate readiness wait
		time.Sleep(2 * time.Second)

		// Respond with 307 to retry the original request
		w.Header().Set("Location", original)
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌐 Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
