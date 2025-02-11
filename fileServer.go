package main

import (
	"log"
	"net/http"
	"os"
)

type FileServer struct {
	addr      string
	uploadDir string
}

// NewFileServer creates a new FileServer instance
func NewFileServer(addr string) *FileServer {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		log.Fatal("UPLOAD_DIR environment variable not set")
	}

	return &FileServer{
		addr:      addr,
		uploadDir: uploadDir,
	}
}

// Run starts the file server
func (s *FileServer) Run() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "/uploads"
		// Restrict the URL to /uploads
		if len(r.URL.Path) < len(prefix) || r.URL.Path[:len(prefix)] != prefix {
			http.NotFound(w, r)
			return
		}

		// Get the file path from the request URL
		filePath := s.uploadDir + r.URL.Path[len(prefix):]

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// Serve the file
		http.ServeFile(w, r, filePath)
	})

	// Start the server on the provided address
	log.Printf("Server running on %s", s.addr)
	err := http.ListenAndServe(s.addr, handler)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}
