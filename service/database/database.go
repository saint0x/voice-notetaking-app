// service/database/database.go

package database

import (
  "database/sql"
  "errors"
  "github.com/yourusername/mare/pkg/database/sqlite"
)

// SaveVoiceNote saves voice note to database.
func SaveVoiceNote(userID int, filePath, transcription string) error {
  // Placeholder for database save logic
  return errors.New("not implemented")
}
