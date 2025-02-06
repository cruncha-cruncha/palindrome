package main

import (
	"net/http"
)

// I like this "server struct" pattern better than closures
type SharedState struct {
    mo *MessageOrchestrator
}

func (ss *SharedState) SaveMessage(w http.ResponseWriter, r *http.Request) {
	// call ss.mo.Add(text)
	// return 201 CreateMessageResponseData

	w.WriteHeader(http.StatusNotImplemented)
}

func (ss *SharedState) GetMessage(w http.ResponseWriter, r *http.Request) {
	// read id from path
	// call ss.mo.Get(id)
	// return 200 GetMessageResponseData or 404 with no body

	w.WriteHeader(http.StatusNotImplemented)
}

func (ss *SharedState) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	// read id from path
	// call ss.mo.Update(id, text)
	// return 200 or 404, neither with a body

	w.WriteHeader(http.StatusNotImplemented)
}

func (ss *SharedState) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// read id from path
	// call ss.mo.Delete(id)
	// return 204 or 404, neither with a body

	w.WriteHeader(http.StatusNotImplemented)
}

func (ss *SharedState) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	// call ss.mo.GetAll()
	// return 200 GetAllMessagesResponseData

	w.WriteHeader(http.StatusNotImplemented)
}
