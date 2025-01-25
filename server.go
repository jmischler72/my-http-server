package main

import (
 "database/sql"   // Package for SQL database interactions
 "fmt"            // Package for formatted I/O
 "log"            // Package for logging
 "net/http"       // Package for HTTP client and server
 "text/template"  // Package for HTML templates

 _ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Todo represents a task in the todo list
type Todo struct {
 ID    int
 Title string
}

// DB is a global variable for the SQLite database connection
var DB *sql.DB

// initDB initializes the SQLite database and creates the todos table if it doesn't exist
func initDB() {
 var err error
 DB, err = sql.Open("sqlite3", "./app.db") // Open a connection to the SQLite database file named app.db
 if err != nil {
  log.Fatal(err) // Log an error and stop the program if the database can't be opened
 }

 // SQL statement to create the todos table if it doesn't exist
 sqlStmt := `
 DROP TABLE IF EXISTS todos;
 CREATE TABLE  todos (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  title TEXT,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
 );`

 _, err = DB.Exec(sqlStmt)
 if err != nil {
  log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt) // Log an error if table creation fails
 }
}

// indexHandler serves the main page and displays all todos
func indexHandler(w http.ResponseWriter, r *http.Request) {
 // Query the database to get all todos
 rows, err := DB.Query("SELECT id, title FROM todos ORDER BY timestamp DESC LIMIT 1;")
 if err != nil { 
  http.Error(w, err.Error(), http.StatusInternalServerError) // Return an HTTP 500 error if the query fails
  return
 }
 defer rows.Close() // Ensure rows are closed after processing

 todos := []Todo{} // Slice to store todos
 for rows.Next() {
  var todo Todo
  if err := rows.Scan(&todo.ID, &todo.Title); err != nil {
   http.Error(w, err.Error(), http.StatusInternalServerError) // Return an HTTP 500 error if scanning fails
   return
  }
  todos = append(todos, todo)
 }

 // Parse and execute the HTML template with the todos data
 tmpl := template.Must(template.New("index").Parse(`
 <!DOCTYPE html>
 <html>
 <head>
  <title>Todo List</title>
 </head>
 <body>
    <h1>jmischler72</h1><ul><li><a href='https://www.youtube.com/@jmischler72'>yt</a></li><li><a href='https://jmischler72.github.io/cv/'>cv</a></li></ul>
    <p>leave a message for the next person</p>
    <form action="/change" method="POST">
   <input type="text" name="title" placeholder="message" required>
   <button type="submit">send</button>
  </form>
  <ul>
   {{range .}}
   <li>{{.Title}}</li>
   {{end}}
  </ul>
 </body>
 </html>
 `))

 tmpl.Execute(w, todos) // Render the template with the list of todos
}

// createHandler handles the creation of a new todo
func createHandler(w http.ResponseWriter, r *http.Request) {
 if r.Method == "POST" {
  title := r.FormValue("title") // Get the title from the form data
  _, err := DB.Exec("INSERT INTO todos(title) VALUES(?)", title) // Insert the new todo into the database
  if err != nil {
   http.Error(w, err.Error(), http.StatusInternalServerError) // Return an HTTP 500 error if insertion fails
   return
  }
  http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the main page after successful creation
 }
}



func main() {
 initDB()         // Initialize the database
 defer DB.Close() // Ensure the database connection is closed when the program exits

 // Route the handlers for each URL path
 http.HandleFunc("/", indexHandler)
 http.HandleFunc("/change", createHandler)

 fmt.Println("Server is running at http://localhost:8080")
 log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server on port 8080
}