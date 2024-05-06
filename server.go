package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/list", handleListFiles) // New endpoint to list files

	logInfo("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logInfo(fmt.Sprintf("Error starting server: %s", err))
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
		logInfo(fmt.Sprintf("Upload failed: %s", err))
		return
	}
	defer uploadedFile.Close()

	uploadPath := "uploads/"
	os.MkdirAll(uploadPath, os.ModePerm)

	targetFilePath := filepath.Join(uploadPath, fileHeader.Filename)
	newFile, err := os.Create(targetFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating the new file: %s", err), http.StatusInternalServerError)
		logInfo(fmt.Sprintf("Error creating file: %s", err))
		return
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, uploadedFile); err != nil {
		http.Error(w, fmt.Sprintf("Error saving the uploaded file: %s", err), http.StatusInternalServerError)
		logInfo(fmt.Sprintf("Error in saving file: %s", err))
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", targetFilePath)
	logInfo(fmt.Sprintf("Uploaded: %s", targetFilePath))
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
		logInfo(fmt.Sprintf("Download failed: %s", err))
		return
	}
	defer targetFile.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+requestedFileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, targetFile); err != nil {
		http.Error(w, fmt.Sprintf("Error sending the requested file: %s", err), http.StatusInternalServerError)
		logInfo(fmt.Sprintf("Error in sending file: %s", err))
		return
	}
	logInfo(fmt.Sprintf("Downloaded: %s", filePath))
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
		logInfo(fmt.Sprintf("Delete failed: %s", err))
		return
	}

	fmt.Fprintf(w, "File deleted successfully: %s", fileNameToDelete)
	logInfo(fmt.Sprintf("Deleted: %s", fileNameToDelete))
}

// New function to list all uploaded files
func handleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadPath := "uploads/"
	files, err := os.ReadDir(uploadPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading the upload directory: %s", err), http.StatusInternalServerError)
		logInfo(fmt.Sprintf("Error listing files: %s", err))
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	fmt.Fprintf(w, "Files: %s", strings.Join(fileList, ", "))
	logInfo(fmt.Sprintf("Listed files"))
}

func logInfo(message string) {
	log.Printf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
}