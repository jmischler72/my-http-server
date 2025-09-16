package main

import (
 "database/sql"   // Package for SQL database interactions
 "encoding/json"  // Package for JSON encoding/decoding
 "fmt"            // Package for formatted I/O
 "log"            // Package for logging
 "net/http"       // Package for HTTP client and server
 "text/template"  // Package for HTML templates

 _ "github.com/mattn/go-sqlite3" // SQLite driver
)

// GridEntry represents a user's entry in the ASCII grid
type GridEntry struct {
 ID       int    `json:"id"`
 X        int    `json:"x"`
 Y        int    `json:"y"`
 Name     string `json:"name"`
 Message  string `json:"message"`
}

// DB is a global variable for the SQLite database connection
var DB *sql.DB

// tmpl is a global variable for the HTML template
var tmpl *template.Template

// initDB initializes the SQLite database and creates the todos table if it doesn't exist
func initDB() {
 var err error
 DB, err = sql.Open("sqlite3", "./db/app.db") // Open a connection to the SQLite database file named app.db
 if err != nil {
  log.Fatal(err) // Log an error and stop the program if the database can't be opened
 }

 // SQL statement to create the grid_entries table for ASCII grid
 gridStmt := `
 CREATE TABLE IF NOT EXISTS grid_entries (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  x INTEGER NOT NULL,
  y INTEGER NOT NULL,
  name TEXT NOT NULL,
  message TEXT NOT NULL,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
 );`

 _, err = DB.Exec(gridStmt)
 if err != nil {
  log.Fatalf("Error creating grid_entries table: %q: %s\n", err, gridStmt)
 }

 // Load HTML template
 tmpl, err = template.ParseFiles("templates/index.html")
 if err != nil {
  log.Fatalf("Error loading template: %v", err)
 }
}

// indexHandler serves the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
 // Execute the HTML template
 tmpl.Execute(w, nil) // Render the template
}



// gridEntriesHandler returns all grid entries as JSON
func gridEntriesHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method != "GET" {
  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  return
 }

 rows, err := DB.Query("SELECT id, x, y, name, message FROM grid_entries ORDER BY timestamp DESC")
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 defer rows.Close()

 entries := []GridEntry{}
 for rows.Next() {
  var entry GridEntry
  if err := rows.Scan(&entry.ID, &entry.X, &entry.Y, &entry.Name, &entry.Message); err != nil {
   http.Error(w, err.Error(), http.StatusInternalServerError)
   return
  }
  entries = append(entries, entry)
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(entries)
}

// gridEntryHandler handles creating a new grid entry
func gridEntryHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method != "POST" {
  http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  return
 }

 var entry GridEntry
 if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Invalid JSON data",
  })
  return
 }

 // Validate input
 if entry.X < 0 || entry.X >= 80 || entry.Y < 0 || entry.Y >= 25 {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Position out of bounds",
  })
  return
 }

 if entry.Name == "" || entry.Message == "" {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Name and message are required",
  })
  return
 }

 // Check if position is already occupied
 var count int
 err := DB.QueryRow("SELECT COUNT(*) FROM grid_entries WHERE x = ? AND y = ?", entry.X, entry.Y).Scan(&count)
 if err != nil {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Database error",
  })
  return
 }

 if count > 0 {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Position already occupied",
  })
  return
 }

 // Insert the new grid entry
 _, err = DB.Exec("INSERT INTO grid_entries(x, y, name, message) VALUES(?, ?, ?, ?)", 
  entry.X, entry.Y, entry.Name, entry.Message)
 if err != nil {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
   "success": false,
   "error":   "Failed to save entry",
  })
  return
 }

 w.Header().Set("Content-Type", "application/json")
 json.NewEncoder(w).Encode(map[string]interface{}{
  "success": true,
 })
}



func main() {
 initDB()         // Initialize the database
 defer DB.Close() // Ensure the database connection is closed when the program exits

 // Route the handlers for each URL path
 http.HandleFunc("/", indexHandler)
 http.HandleFunc("/grid-entries", gridEntriesHandler)
 http.HandleFunc("/grid-entry", gridEntryHandler)
 
 // Serve static files (CSS, JS)
 http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

 fmt.Println("Server is running at http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server on port 8080
}