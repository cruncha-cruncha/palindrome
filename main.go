package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	ss := &SharedState{}

	r.Methods("POST").Path("/messages").HandlerFunc(ss.SaveMessage)
	r.Methods("GET").Path("/messages").HandlerFunc(ss.GetAllMessages)
	r.Methods("GET").Path("/messages/{mid}").HandlerFunc(ss.GetMessage)
	r.Methods("PUT").Path("/messages/{mid}").HandlerFunc(ss.UpdateMessage) // not PATCH
	r.Methods("DELETE").Path("/messages/{mid}").HandlerFunc(ss.DeleteMessage)

	http.ListenAndServe(":8090", r)
}