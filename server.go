package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/delete", handleDelete)

	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadedFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the uploaded file: %s", err), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	uploadPath := "uploads/"
	os.MkdirAll(uploadPath, os.ModePerm)

	targetFilePath := filepath.Join(uploadPath, fileHeader.Filename)
	newFile, err := os.Create(targetFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating the new file: %s", err), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, uploadedFile); err != nil {
		http.Error(w, fmt.Sprintf("Error saving the uploaded file: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", targetFilePath)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestedFileName := r.URL.Query().Get("file")
	if requestedFileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("uploads", requestedFileName)
	targetFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening the requested file: %s", err), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+requestedFileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, targetFile); err != nil {
		http.Error(w, fmt.Sprintf("Error sending the requested file: %s", err), http.StatusInternalServerError)
		return
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileNameToDelete := r.URL.Query().Get("file")
	if fileNameToDelete == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	targetFilePath := filepath.Join("uploads", fileNameToDelete)
	if err := os.Remove(targetFilePath); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting the file: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File deleted successfully: %s", fileNameToDelete)
}