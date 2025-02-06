package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ParseIdFromPath(r *http.Request) (int, error) {
	params := mux.Vars(r)
	str_id, ok := params["id"]
	if !ok {
		return 0, errors.New("id not found in path")
	}

	id, err := strconv.Atoi(str_id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// conveniece functions for converting between types

func NewGetMessageResponseDataFromMessage(m *Message) GetMessageResponseData {
	return GetMessageResponseData{
		Text: m.text,
		IsPalindrome: PalindromeStatusToBoolPointer(m.isPalindrome),
	}
}

func NewGetMessagesResponseItemFromMessage(m *Message) GetAllMessagesResponseItem {
	return GetAllMessagesResponseItem{
		ID: m.id,
		Text: m.text,
		IsPalindrome: PalindromeStatusToBoolPointer(m.isPalindrome),
	}
}
