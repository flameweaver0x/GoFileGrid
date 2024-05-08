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
    http.HandleFunc("/upload", handleUpload) // Adapted for multiple files
    http.HandleFunc("/download", handleDownload)
    http.HandleFunc("/delete", handleDelete)
    http.HandleFunc("/list", handleListFiles)

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

    // Parse the multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, fmt.Sprintf("Error parsing multipart form: %s", err), http.StatusInternalServerError)
        logInfo(fmt.Sprintf("Error parsing multipart form: %s", err))
        return
    }

    uploadPath := "uploads/"
    os.MkdirAll(uploadPath, os.ModePerm)

    files := r.MultipartForm.File["files"] // Note the field name is now expected to be "files" and a slice
    var uploadedFiles []string
    for _, fileHeader := range files {
        uploadedFile, err := fileHeader.Open()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error retrieving the uploaded file: %s", err), http.StatusInternalServerError)
            logInfo(fmt.Sprintf("Upload failed: %s", err))
            continue
        }
        defer uploadedFile.Close()

        targetFilePath := filepath.Join(uploadPath, fileHeader.Filename)
        newFile, err := os.Create(targetFilePath)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error creating the new file: %s", err), http.StatusInternalServerError)
            logInfo(fmt.Sprintf("Error creating file: %s", err))
            continue
        }
        defer newFile.Close()

        if _, err := io.Copy(newFile, uploadedFile); err != nil {
            http.Error(w, fmt.Sprintf("Error saving the uploaded file: %s", err), http.StatusInternalServerError)
            logInfo(fmt.Sprintf("Error in saving file: %s", err))
            return
        }
        logInfo(fmt.Sprintf("Uploaded: %s", targetFilePath))
        uploadedFiles = append(uploadedFiles, fileHeader.Filename)
    }

    fmt.Fprintf(w, "Files uploaded successfully: %s", strings.Join(uploadedFiles, ", "))
}

func logInfo(message string) {
    log.Printf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
}