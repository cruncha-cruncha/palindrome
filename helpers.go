package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// CalculateHash returns the SHA-256 hash of some given text.
func CalculateHash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// ParseIdFromPath extracts the "id" parameter from the request path. It uses 
// the gorilla/mux package. It returns 0 and an error if the "id" parameter is
// not found or if it is not a valid integer.
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

// BinarySearch performs a binary search on a slice of any type, assuming that
// it's already sorted. The selector function is used to determine an elements
// value for the purpose of comparison. So every element E has a an associated
// integer I.
// 
// BinarySearch returns the index of some element E having I == target. It could
// be the first index, the last, or one inbetween (there's no guarantee). If no
// suitable element E is found, it returns the index where an element with I ==
// target should be inserted to maintain the sorted order.
//
// It's useful when adding elements to an already sorted slice, or when building
// a sorted slice one element at a time.
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

// PWorkKeyFromMsg is a convenience function which creates a PalindromeWorkKey
// from a Message.
func PWorkKeyFromMsg(msg Message) PalindromeWorkKey {
	return PalindromeWorkKey{
		hash:      msg.hash,
		messageId: msg.id,
	}
}
