package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/delete", deleteHandler)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %s", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	dir := "uploads/"
	os.MkdirAll(dir, os.ModePerm) 

	filePath := filepath.Join(dir, header.Filename)
	newFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating the file: %s", err), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, file); err != nil {
		http.Error(w, fmt.Sprintf("Error saving the file: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", filePath)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening the file: %s", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, fmt.Sprintf("Error sending the file: %s", err), http.StatusInternalServerError)
		return
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", fileName)
	if err := os.Remove(filePath); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting the file: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File deleted successfully: %s", fileName)
}