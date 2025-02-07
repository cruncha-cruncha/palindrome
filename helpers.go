package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CalculateHash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

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

func BinarySearch[T any](arr []T, selector func(*T) int, target int) int {
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

func PWorkKeyFromMsg(msg Message) PalindromeWorkKey {
	return PalindromeWorkKey{
		hash:      msg.hash,
		messageId: msg.id,
	}
}
