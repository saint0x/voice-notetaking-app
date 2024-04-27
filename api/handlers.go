// api/handlers.go

package api

import (
  "encoding/json"
  "io"
  "net/http"
  "github.com/yourusername/mare/service/speechtotext"
)

// UploadVoiceNoteHandler handles voice note uploads.
func UploadVoiceNoteHandler(w http.ResponseWriter, r *http.Request) {
  // Parse multipart form data
  err := r.ParseMultipartForm(10 << 20) // 10 MB
  if err != nil {
    http.Error(w, "Failed to parse form data", http.StatusBadRequest)
    return
  }

  // Get audio file from form data
  file, _, err := r.FormFile("audio")
  if err != nil {
    http.Error(w, "Failed to read audio file", http.StatusBadRequest)
    return
  }
  defer file.Close()

  // Read file content
  fileBytes, err := io.ReadAll(file)
  if err != nil {
    http.Error(w, "Failed to read audio file content", http.StatusInternalServerError)
    return
  }

  // Convert audio to text using AssemblyAI
  transcription, err := speechtotext.UploadToAssemblyAI(fileBytes)
  if err != nil {
    http.Error(w, "Failed to convert speech to text", http.StatusInternalServerError)
    return
  }

  // Return transcription as JSON response
  response := map[string]string{
    "transcription": transcription,
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(response)
}
