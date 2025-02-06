package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

func (ss *SharedState) SaveMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload CreateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err, _, _ := ss.mo.Add(payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateMessageResponseData{ID: id})
}

func (ss *SharedState) GetMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewGetMessageResponseDataFromMessage(&message))
}

func (ss *SharedState) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload UpdateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, _, _ = ss.mo.Update(id, payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ss *SharedState) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ss.mo.Delete(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ss *SharedState) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	data := GetAllMessagesResponseData{
		Messages: []GetAllMessagesResponseItem{},
	}

	messages, err := ss.mo.GetAll()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, m := range messages {
		// sort while we insert
		insertIndex := binarySearch(data.Messages, func(m *GetAllMessagesResponseItem) int { return m.ID }, m.id)
		data.Messages = slices.Insert(data.Messages, insertIndex, NewGetMessagesResponseItemFromMessage(&m))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
