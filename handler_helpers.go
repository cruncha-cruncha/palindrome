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


func binarySearch[T any](arr []T, selector func(*T) int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := (left + right) / 2
		midValue := selector(&arr[mid])
		if midValue == target {
			return mid
		} else if midValue < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return left
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
