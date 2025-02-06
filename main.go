package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	ss := NewSharedState()

	r.Methods("POST").Path("/messages").HandlerFunc(ss.SaveMessage)
	r.Methods("GET").Path("/messages").HandlerFunc(ss.GetAllMessages)
	r.Methods("GET").Path("/messages/{id}").HandlerFunc(ss.GetMessage)
	r.Methods("PUT").Path("/messages/{id}").HandlerFunc(ss.UpdateMessage) // not PATCH
	r.Methods("DELETE").Path("/messages/{id}").HandlerFunc(ss.DeleteMessage)

	http.ListenAndServe(":8090", r)
}