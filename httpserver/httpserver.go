package httpserver // github.com/microbuilder/elfquery/httpserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microbuilder/elfquery/elf2sql"
)

// REST API catch all handler
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "endpoint not found"}`))
}

// Root page handler
func home(w http.ResponseWriter, r *http.Request) {
	// Query the database
	query := "SELECT Name, Type, Binding, Visibility, Section, printf('0x%X', Value) AS Address, Size FROM symbols ORDER BY Size DESC LIMIT 50"
	s, e := elf2sql.RunQuery(query, elf2sql.DFHtml)
	if e != nil {
		w.Write([]byte("Invalid query\n"))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Load the template
	tmpl, e := template.ParseFiles("templates/query.html")
	if e != nil {
		fmt.Printf("Unable to load template file.\n")
		return
	}

	// Enable searching in table results via tbody id
	s = strings.Replace(s, "<tbody>", "<tbody id=\"restable\">", 1)

	// Inject results
	data := struct {
		PageTitle string
		SQLQuery  string
		Results   template.HTML
	}{
		PageTitle: "Query Results",
		SQLQuery:  query,
		Results:   template.HTML(s),
	}
	tmpl.Execute(w, data)
}

// Start the HTTP Server
func Start(port int16) {
	r := mux.NewRouter()

	// Setup the REST API subrouter
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", notFound)

	// Handle standard requests. Routes are tested in the order they are added,
	// so these will only be handled if they don't match anything above.
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js/"))))
	r.HandleFunc("/", home)

	fmt.Println("Starting HTTP server on port http://localhost:" + strconv.Itoa(int(port)))
	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
