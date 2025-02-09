// Package main contains everything in this project
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Main sets up routing, shared state, and starts the server. It listens on port
// 8090 by default, overridden by the PORT environment variable.
func main() {
	r := mux.NewRouter()
	ss := NewSharedState()

	r.Methods("POST").Path("/messages").HandlerFunc(ss.CreateMessage)
	r.Methods("GET").Path("/messages").HandlerFunc(ss.GetAllMessages)
	r.Methods("DELETE").Path("/messages").HandlerFunc(ss.DeleteAllMessages)
	r.Methods("GET").Path("/messages/{id}").HandlerFunc(ss.GetMessage)
	r.Methods("PUT").Path("/messages/{id}").HandlerFunc(ss.UpdateMessage) // not PATCH, as we're effectively replacing the whole message
	r.Methods("DELETE").Path("/messages/{id}").HandlerFunc(ss.DeleteMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("Listening on port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	log.Fatal(err)
}
