package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

// SaveMessage expects a JSON payload with a "text" field and returns 200 with
// a JSON response, which has an "id" field (a positive integer).
func (ss *SharedState) SaveMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload CreateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, err := ss.mo.Add(payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, _, err = ss.po.Add(msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateMessageResponseData{ID: msg.id})
}

// GetMessage expects an ID in the path and returns a JSON response with two
// fields: "text" and "is_palindrome". The "is_palindrome" field is a boolean
// but can be null. It will return 404 if the message is not found.
func (ss *SharedState) GetMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	workKey := PWorkKeyFromMsg(msg)
	found, result, _, err := ss.po.Poll(workKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetMessageResponseData{
		Text:         msg.text,
		IsPalindrome: PStatusToBoolPointer(result.isPalindrome),
	})
}

// UpdateMessage expects an ID in the path as well as a JSON payload with a
// "text" field. It will return 404 if the message to be updated is not found,
// otherwise it will return 200, no body.
func (ss *SharedState) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload UpdateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldMsg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newMsg, err := ss.mo.Update(id, payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, _, err = ss.po.Add(newMsg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	oldWorkKey := PWorkKeyFromMsg(oldMsg)
	err = ss.po.Remove(oldWorkKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteMessage expects an ID in the path. It will return 404 if the message
// to be deleted doesn't exist, otherwise it will return 204, no body.
func (ss *SharedState) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = ss.mo.Delete(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	workKey := PWorkKeyFromMsg(msg)
	err = ss.po.Remove(workKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllMessages returns a JSON response with a 'messages' field, which is an
// array of objects with 'id', 'text', and 'is_palindrome' fields. The array
// is sorted by 'id' in ascending order.
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
		workKey := PWorkKeyFromMsg(m)
		found, result, _, err := ss.po.Poll(workKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !found {
			result = PalindromeWorkStatus{isPalindrome: P_UNKNOWN}
		}

		// sort while we insert
		insertIndex := BinarySearch(data.Messages, func(m *GetAllMessagesResponseItem) int { return m.ID }, m.id)
		data.Messages = slices.Insert(data.Messages, insertIndex, GetAllMessagesResponseItem{
			ID:           m.id,
			Text:         m.text,
			IsPalindrome: PStatusToBoolPointer(result.isPalindrome),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// DeleteAllMessages deletes all messages and returns 204, no body.
func (ss *SharedState) DeleteAllMessages(w http.ResponseWriter, r *http.Request) {
	err := ss.mo.DeleteAll()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ss.po.Clear()

	w.WriteHeader(http.StatusNoContent)
}
