package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

// CreateMessage expects a JSON payload with a "text" field and returns 200 with
// a JSON response, which has an "id" field (a positive integer).
func (ss *SharedState) CreateMessage(w http.ResponseWriter, r *http.Request) {
	// verify payload (need some text)
	decoder := json.NewDecoder(r.Body)
	var payload CreateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create the message
	msg, err := ss.mo.Add(payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// kick off the palindrome work
	_, _, _, err = ss.po.Add(msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// respond with message id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateMessageResponseData{ID: msg.id})
}

// GetMessage expects an ID in the path and returns a JSON response with two
// fields: "text" and "is_palindrome". The "is_palindrome" field is a boolean
// but can be null. It will return 404 if the message is not found.
func (ss *SharedState) GetMessage(w http.ResponseWriter, r *http.Request) {
	// get the message id we're looking for
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the message, return 404 if not found
	msg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// get the palindrome work result corresponding to the message
	workKey := PWorkKeyFromMsg(msg)
	found, result, _, err := ss.po.Poll(workKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		// This can happen on Add or Update: adding new work requires a message
		// id, which doesn't exist until the message is created, so it's 
		// possible for a message to exist with no corresponding palindrome work
		// (for a brief moment in time).
		// 
		// It shouldn't happen on Delete or DeleteAll: messages are deleted
		// before their work. But it's also a possible bug. In any case, there's
		// no harm in inserting more work (duplicate work is handled / ignored).

		// Kick off more work, so next time we'll have a result.
		ss.po.Add(msg)

		// Can safely return P_UNKNOWN, even though it could be annoying for the
		// user.
		result = PWResult{isPalindrome: P_UNKNOWN}
	}

	// respond with the message text and palindrome status
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
	// get the message id we want to update
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify payload (need some text)
	decoder := json.NewDecoder(r.Body)
	var payload UpdateMessageRequestData
	if err := decoder.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify that we're updating an existing message
	oldMsg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// update the message
	newMsg, err := ss.mo.Update(id, payload.Text)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// kick off palindrome work for the new message
	_, _, _, err = ss.po.Add(newMsg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// cancel palindrome work for the old message
	oldWorkKey := PWorkKeyFromMsg(oldMsg)
	err = ss.po.Remove(oldWorkKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// respond
	w.WriteHeader(http.StatusOK)
}

// DeleteMessage expects an ID in the path. It will return 404 if the message
// to be deleted doesn't exist, otherwise it will return 204, no body.
func (ss *SharedState) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// get the message id we want to delete
	id, err := ParseIdFromPath(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// verify that the message exists
	msg, found, err := ss.mo.Get(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// delete the message
	err = ss.mo.Delete(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// cancel the corresponding palindrome work
	workKey := PWorkKeyFromMsg(msg)
	err = ss.po.Remove(workKey)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// respond
	w.WriteHeader(http.StatusNoContent)
}

// GetAllMessages returns a JSON response with a 'messages' field, which is an
// array of objects with 'id', 'text', and 'is_palindrome' fields. The array
// is sorted by 'id' in ascending order.
func (ss *SharedState) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	// no message id or payload to parse

	// this will be our response data
	data := GetAllMessagesResponseData{
		Messages: []GetAllMessagesResponseItem{},
	}

	// get all messages
	messages, err := ss.mo.GetAll()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// for each message, get the corresponding palindrome work, then format and
	// add to the response data
	for _, m := range messages {
		workKey := PWorkKeyFromMsg(m)
		found, result, _, err := ss.po.Poll(workKey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !found {
			// This should never happen, but we can handle it. See GetMessage
			// for more details.
			result = PWResult{isPalindrome: P_UNKNOWN}
			ss.po.Add(m)
		}

		// Sort the response while we insert. Messages will end up in ascending
		// order by ID.
		insertIndex := BinarySearch(data.Messages, func(m *GetAllMessagesResponseItem) int { return m.ID }, m.id)
		data.Messages = slices.Insert(data.Messages, insertIndex, GetAllMessagesResponseItem{
			ID:           m.id,
			Text:         m.text,
			IsPalindrome: PStatusToBoolPointer(result.isPalindrome),
		})
	}

	// respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// DeleteAllMessages deletes all messages and returns 204, no body.
func (ss *SharedState) DeleteAllMessages(w http.ResponseWriter, r *http.Request) {
	// no message id or payload to parse

	// delete all messages
	err := ss.mo.DeleteAll()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// cancel all palindrome work
	err = ss.po.Clear()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// respond
	w.WriteHeader(http.StatusNoContent)
}
