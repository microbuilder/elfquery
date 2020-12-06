package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// REST API catch all handler
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "endpoint not found"}`))
}

// Root page handler
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("HTTP server\n"))
}

// Start the HTTP Server
func Start(port int16) {
	r := mux.NewRouter()

	// Setup the REST API subrouter
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", notFound)

	// Handle standard requests. Routes are tested in the order they are added,
	// so these will only be handled if they don't match anything above.
	r.HandleFunc("/", home)

	// Make sure the server key and certificate exist
	if !fileExists("SERVER.key") || !fileExists("SERVER.crt") {
		log.Fatal("SERVER.crt and/or SSERVER.key not found. See README.md.")
	}

	fmt.Println("Starting HTTPS server on port https://localhost:" + strconv.Itoa(int(port)))
	err := http.ListenAndServeTLS(":"+strconv.Itoa(int(port)), "SERVER.crt", "SERVER.key", r)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
