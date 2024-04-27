// pkg/database/sqlite/sqlite.go

package sqlite

import (
  "database/sql"
  "log"
  "time"

  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize initializes the SQLite database and tables.
func Initialize(dbPath string) error {
  var err error

  db, err = sql.Open("sqlite3", dbPath)
  if err != nil {
    return err
  }

  // Create tables if they don't exist
  err = InitializeTables()
  if err != nil {
    return err
  }

  log.Println("SQLite database initialized successfully")

  return nil
}

// InitializeTables initializes database tables if they don't exist.
func InitializeTables() error {
  _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT UNIQUE NOT NULL,
      password_hash TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS voice_notes (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER,
      file_path TEXT NOT NULL,
      transcription TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS topics (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT UNIQUE NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS tags (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT UNIQUE NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS voice_note_topics (
      voice_note_id INTEGER,
      topic_id INTEGER,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (voice_note_id, topic_id),
      FOREIGN KEY (voice_note_id) REFERENCES voice_notes(id),
      FOREIGN KEY (topic_id) REFERENCES topics(id)
    );

    CREATE TABLE IF NOT EXISTS voice_note_tags (
      voice_note_id INTEGER,
      tag_id INTEGER,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (voice_note_id, tag_id),
      FOREIGN KEY (voice_note_id) REFERENCES voice_notes(id),
      FOREIGN KEY (tag_id) REFERENCES tags(id)
    );
  `)
  if err != nil {
    return err
  }

  log.Println("Database tables initialized successfully")

  return nil
}

// InsertVoiceNote inserts a new voice note into the database.
func InsertVoiceNote(userID int, filePath, transcription string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO voice_notes (user_id, file_path, transcription) 
    VALUES (?, ?, ?)
  `, userID, filePath, transcription)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// InsertTopic inserts a new topic into the database.
func InsertTopic(name string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO topics (name) 
    VALUES (?)
  `, name)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// InsertTag inserts a new tag into the database.
func InsertTag(name string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO tags (name) 
    VALUES (?)
  `, name)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// InsertVoiceNoteTopic inserts a connection between a voice note and a topic into the database.
func InsertVoiceNoteTopic(voiceNoteID, topicID int64) error {
  _, err := db.Exec(`
    INSERT INTO voice_note_topics (voice_note_id, topic_id) 
    VALUES (?, ?)
  `, voiceNoteID, topicID)
  if err != nil {
    return err
  }

  return nil
}

// InsertVoiceNoteTag inserts a connection between a voice note and a tag into the database.
func InsertVoiceNoteTag(voiceNoteID, tagID int64) error {
  _, err := db.Exec(`
    INSERT INTO voice_note_tags (voice_note_id, tag_id) 
    VALUES (?, ?)
  `, voiceNoteID, tagID)
  if err != nil {
    return err
  }

  return nil
}

// Close closes the SQLite database connection.
func Close() error {
  if db != nil {
    return db.Close()
  }
  return nil
}
