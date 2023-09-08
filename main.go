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
	// "time"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("Server started...")

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

	handler2 := func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(1 * time.Second)
		// log.Print("HTMX request received")
		// log.Print(r.Header.Get("HX-Request")) // Boolean
		title := r.PostFormValue("film-title")
		director := r.PostFormValue("film-director")
		fmt.Print(title)
		fmt.Print(director)
		// NOTE: Typically will add this data to a db
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
