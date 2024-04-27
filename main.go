package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "path/filepath"

  "service/tagging/tagging.go"
      
"service/summarization/summarization"
  "service/timelogging/timelogging"
  "service/speechtotext/speechtotext.go"
  "service/insight/insight.go"
)

func main() {
  // Initialize SQLite database
  err := sqlite.Initialize("pkg/database/sqlite/sqlite.go")
  if err != nil {
    log.Fatalf("Failed to initialize SQLite database: %v", err)
  }
  defer sqlite.Close()

  // HTTP handler to upload voice note
  http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form data
    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
      log.Printf("Failed to parse form data: %v", err)
      http.Error(w, "Failed to parse form data", http.StatusBadRequest)
      return
    }

    // Get audio file from form data
    file, _, err := r.FormFile("audio")
    if err != nil {
      log.Printf("Failed to read audio file: %v", err)
      http.Error(w, "Failed to read audio file", http.StatusBadRequest)
      return
    }
    defer file.Close()

    // Read file content
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
      log.Printf("Failed to read file content: %v", err)
      http.Error(w, "Failed to read file content", http.StatusInternalServerError)
      return
    }

    // Convert audio to text using speech-to-text service
    transcription, err := speechtotext.UploadToAssemblyAI(fileBytes)
    if err != nil {
      log.Printf("Failed to transcribe audio: %v", err)
      http.Error(w, "Failed to transcribe audio", http.StatusInternalServerError)
      return
    }
    log.Println("Transcription:", transcription)

    // Summarize the transcription
    summary, err := summarization.SummarizeText(transcription)
    if err != nil {
      log.Printf("Failed to summarize text: %v", err)
      http.Error(w, "Failed to summarize text", http.StatusInternalServerError)
      return
    }
    log.Println("Summary:", summary)

    // Tag the transcription
    tags, err := tagging.TagText(transcription)
    if err != nil {
      log.Printf("Failed to tag text: %v", err)
      http.Error(w, "Failed to tag text", http.StatusInternalServerError)
      return
    }
    log.Println("Tags:", tags)

    // Generate insights from the transcription
    groupedNotes := []string{transcription} // For simplicity, using the entire transcription as one note
    insightText, err := insight.GenerateInsight(groupedNotes)
    if err != nil {
      log.Printf("Failed to generate insight: %v", err)
      http.Error(w, "Failed to generate insight", http.StatusInternalServerError)
      return
    }
    log.Println("Insight:", insightText)

    // Store results in SQLite database
    // Here you can call SQLite database functions to insert the data
    // ...

    // Send success response
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Processing completed successfully!")
  })

  // Start HTTP server
  log.Println("Server is running on port 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
