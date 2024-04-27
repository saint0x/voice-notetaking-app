// service/speechtotext/speechtotext.go

package speechtotext

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "io"
  "net/http"
)

// UploadToAssemblyAI uploads audio file to AssemblyAI and returns transcription.
func UploadToAssemblyAI(fileBytes []byte) (string, error) {
  // Encode file to base64
  encodedFile := base64.StdEncoding.EncodeToString(fileBytes)

  // Prepare request body
  requestBody := map[string]string{
    "data": "data:application/octet-stream;base64," + encodedFile,
  }

  // Marshal request body to JSON
  jsonBody, err := json.Marshal(requestBody)
  if err != nil {
    return "", err
  }

  // Create HTTP client
  client := &http.Client{}

  // Create POST request
  req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", bytes.NewBuffer(jsonBody))
  if err != nil {
    return "", err
  }

  // Set headers
  req.Header.Set("Authorization", "Bearer a367f96876fe4320afb97af3d989f2b7")
  req.Header.Set("Content-Type", "application/json")

  // Send request
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  // Read response body
  respBody, err := io.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }

  // Parse JSON response
  var jsonResponse map[string]interface{}
  err = json.Unmarshal(respBody, &jsonResponse)
  if err != nil {
    return "", err
  }

  // Check if transcription URL exists in response
  transcriptURL, ok := jsonResponse["upload_url"].(string)
  if !ok {
    return "", errors.New("failed to get upload URL from AssemblyAI response")
  }

  // Call AssemblyAI transcription endpoint
  transcription, err := TranscribeFromAssemblyAI(transcriptURL)
  if err != nil {
    return "", err
  }

  return transcription, nil
}

// TranscribeFromAssemblyAI transcribes audio file from AssemblyAI URL.
func TranscribeFromAssemblyAI(url string) (string, error) {
  // Prepare POST request to AssemblyAI transcription endpoint
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    return "", err
  }

  // Set headers
  req.Header.Set("Authorization", "Bearer a367f96876fe4320afb97af3d989f2b7")

  // Create HTTP client
  client := &http.Client{}

  // Send request
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  // Read response body
  respBody, err := io.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }

  // Parse JSON response
  var jsonResponse map[string]interface{}
  err = json.Unmarshal(respBody, &jsonResponse)
  if err != nil {
    return "", err
  }

  // Get transcription from response
  transcription, ok := jsonResponse["text"].(string)
  if !ok {
    return "", errors.New("failed to get transcription from AssemblyAI response")
  }

  return transcription, nil
}
