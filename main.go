// NOTE:
// - Pass data (usually from db) context from server to template, which then
// gets rendered on the client.
// - Extract POST data from Request.PostFormValue("name-of-input")
// - Default is hx-swap="innerHTML"
//   REF: https://htmx.org/docs/#swapping
// - Can use Template Fragments {{ block }} for small dynamic chunks of HTML,
//   rather than creating a new template file (.html)

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	// "time"
	"context"
	"database/sql"

	// "sync"

	"github.com/joho/godotenv"

	_ "github.com/libsql/libsql-client-go/libsql"
	// _ "modernc.org/sqlite" // for local SQLite files
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("Server started...")

	// Environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	godotenv.Load()

	// Build the DB_URL string from env vars
	// REF: https://github.com/libsql/libsql-client-go/#readme
	// "libsql://smooth-madripoor.turso.io?authToken=[db-auth-token]"
	var sb strings.Builder
	sb.WriteString(os.Getenv("DB_URL"))
	sb.WriteString("?authToken=")
	sb.WriteString(os.Getenv("DB_AUTH_TOKEN"))
	dbURL := sb.String()
	fmt.Println("DB URL is: %v", dbURL)
	if dbURL == "" {
		log.Fatal("DB URL is not found in environment")
	}

	// Open a connection Turso db
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open db %s: %s", dbURL, err)
		os.Exit(1)
	}
	defer db.Close()
	ctx := context.Background()

	// Check if db is accessible (not connected just yet!) - the actual
	// db connection is deferred until a query is made
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	// Try to execute SQL statement to DB
	// REF: https://github.com/libsql/libsql-client-go/blob/42289d60a0305c4c51d33f2e4e6e3b57129881f5/examples/sql/counter/main.go#L14C1-L21C2
	// execStatementNamedArgs := "INSERT INTO movies(id, title, director) VALUES (:id, :title, :director)"
	stmt := "SELECT * FROM movies"
	// res, err := db.ExecContext(ctx, stmt)
	res, err := db.Exec(stmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute statement %s: %s", stmt, err)
		os.Exit(1)
	}
	fmt.Println("Result: ", res)

	// Q: Should I use Query() or Exec()?
	// res, err := db.Query(stmt)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to execute query %s: %s", stmt, err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Result: ", res)

	// GetMovies handler
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		// Q: template.Must() seems similar to Result<()>
		t := template.Must(template.ParseFiles("index.html"))
		// NOTE: map[k:string]v:[]Film -- key is str, val is arr of Film
		// Writing our data to the Template
		films := map[string][]Film{
			"Films": {
				{Title: "The Professional", Director: "Nobody"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}

		t.Execute(w, films)

		// Displaying our data inside Template
		// REF: https://golangforall.com/en/post/templates.html
	}

	// AddMovie Handler
	handler2 := func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(1 * time.Second)
		// log.Print("HTMX request received")
		// log.Print(r.Header.Get("HX-Request")) // Boolean
		title := r.PostFormValue("film-title")
		director := r.PostFormValue("film-director")
		fmt.Print(title)
		fmt.Print(director)
		// NOTE: Typically will add this data to a db
		// U: Let's add record to Turso db
		execStatementNamedArgs := "INSERT INTO movies(id, title, director) VALUES (:id, :title, :director)"
		db.Exec(execStatementNamedArgs, sql.Named("id", `hex(randomblob(8)`), sql.Named("title", title), sql.Named("director", director))

		// Let's return an HTML string
		// U: Rather than hardcoding HTML, better to use a Fragment
		// htmlStr := fmt.Sprintf("<li>%s - %s</li>", title, director)
		// // Create a new Template
		// t, _ := template.New("t").Parse(htmlStr)
		// // NOTE: Don't have any data -- simply rendering <li> element
		// t.Execute(w, nil)
		// U: Let's use a Fragment instead
		t := template.Must(template.ParseFiles("index.html"))
		t.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})

	}

	http.HandleFunc("/", handler1)
	http.HandleFunc("/add-film/", handler2)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
