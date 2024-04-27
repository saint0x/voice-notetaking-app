package sqlite

import (
  "database/sql"
  "log"
  "os"
  "path/filepath"

  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize initializes the SQLite database and tables.
func Initialize(dbPath string) error {
  var err error

  // Create the directory to store the database file if it doesn't exist
  dir := filepath.Dir(dbPath)
  err = os.MkdirAll(dir, os.ModePerm)
  if err != nil {
    return err
  }

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
    CREATE TABLE IF NOT EXISTS nodes (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      text TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS edges (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      source_id INTEGER,
      target_id INTEGER,
      weight FLOAT,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (source_id) REFERENCES nodes(id),
      FOREIGN KEY (target_id) REFERENCES nodes(id)
    );

    CREATE TABLE IF NOT EXISTS vertices (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      node_id INTEGER,
      target_id INTEGER,
      concept TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (node_id) REFERENCES nodes(id),
      FOREIGN KEY (target_id) REFERENCES nodes(id)
    );

    CREATE TABLE IF NOT EXISTS recordings (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER,
      file_path TEXT NOT NULL,
      transcription TEXT NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS concepts (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT UNIQUE NOT NULL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
  `)
  if err != nil {
    return err
  }

  log.Println("Database tables initialized successfully")

  return nil
}

// InsertNode inserts a new node into the database and returns its ID.
func InsertNode(text string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO nodes (text) 
    VALUES (?)
  `, text)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// GetNodeByID retrieves a node from the database by its ID.
func GetNodeByID(id int64) (string, error) {
  var text string
  err := db.QueryRow(`
    SELECT text FROM nodes WHERE id = ?
  `, id).Scan(&text)
  if err != nil {
    return "", err
  }

  return text, nil
}

// InsertEdge inserts a new edge into the database and returns its ID.
func InsertEdge(sourceID, targetID int64, weight float64) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO edges (source_id, target_id, weight) 
    VALUES (?, ?, ?)
  `, sourceID, targetID, weight)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// GetEdgeByID retrieves an edge from the database by its ID.
func GetEdgeByID(id int64) (int64, int64, float64, error) {
  var sourceID, targetID int64
  var weight float64
  err := db.QueryRow(`
    SELECT source_id, target_id, weight FROM edges WHERE id = ?
  `, id).Scan(&sourceID, &targetID, &weight)
  if err != nil {
    return 0, 0, 0, err
  }

  return sourceID, targetID, weight, nil
}

// InsertVertex inserts a new vertex into the database and returns its ID.
func InsertVertex(nodeID, targetID int64, concept string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO vertices (node_id, target_id, concept) 
    VALUES (?, ?, ?)
  `, nodeID, targetID, concept)
  if err != nil {
    return 0, err
  }

  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  return id, nil
}

// GetVertexByID retrieves a vertex from the database by its ID.
func GetVertexByID(id int64) (int64, int64, string, error) {
  var nodeID, targetID int64
  var concept string
  err := db.QueryRow(`
    SELECT node_id, target_id, concept FROM vertices WHERE id = ?
  `, id).Scan(&nodeID, &targetID, &concept)
  if err != nil {
    return 0, 0, "", err
  }

  return nodeID, targetID, concept, nil
}

// InsertRecording inserts a new recording into the database and returns its ID.
func InsertRecording(userID int64, filePath, transcription string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO recordings (user_id, file_path, transcription) 
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

// GetRecordingByID retrieves a recording from the database by its ID.
func GetRecordingByID(id int64) (int64, string, string, error) {
  var userID int64
  var filePath, transcription string
  err := db.QueryRow(`
    SELECT user_id, file_path, transcription FROM recordings WHERE id = ?
  `, id).Scan(&userID, &filePath, &transcription)
  if err != nil {
    return 0, "", "", err
  }

  return userID, filePath, transcription, nil
}

// InsertConcept inserts a new concept into the database and returns its ID.
func InsertConcept(name string) (int64, error) {
  result, err := db.Exec(`
    INSERT INTO concepts (name) 
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

// GetConceptByID retrieves a concept from the database by its ID.
func GetConceptByID(id int64) (string, error) {
  var name string
  err := db.QueryRow(`
    SELECT name FROM concepts WHERE id = ?
  `, id).Scan(&name)
  if err != nil {
    return "", err
  }

  return name, nil
}

// Close closes the SQLite database connection.
func Close() error {
  if db != nil {
    return db.Close()
  }
  return nil
}
