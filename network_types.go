package main

// All incoming and outgoing payloads are JSON. This file contains all the types
// that are converted directly to/from JSON by any handler.

// ---- Request Types ----

// CreateMessageRequestData is used when creating a new message. It has a single
// "text" field.
type CreateMessageRequestData struct {
	Text string `json:"text"`
}

// UpdateMessageRequestData is used when updating an existing message. It has
// a single "text" field, exactly the same as CreateMessageRequestData, but it's
// a separate type for clarity and future-proofing.
type UpdateMessageRequestData struct {
	Text string `json:"text"`
}

// ---- Response Types ----

// CreateMessageResponseData is returned after a new message is created, with
// its new unique id (an integer).
type CreateMessageResponseData struct {
	ID int `json:"id"`
}

// GetMessageResponseData is returned when a message is successfully retrieved.
// It has two fields: "text" and "is_palindrome".
type GetMessageResponseData struct {
	Text         string `json:"text"`
	// IsPalindrome can be null, which means the text is empty, or the server
	// is still calculating whether or not the text is a palindrome. This is
	// trinary logic. In actual production code, I would use an explicit status
	// field instead, this boolean pointer is just for fun.
	IsPalindrome *bool  `json:"is_palindrome"` 
}

// GetAllMessagesResponseData is returned from a request to get all messages. It
// has a single field, "messages", which is an array of
// GetAllMessagesResponseItem.
type GetAllMessagesResponseData struct {
	Messages []GetAllMessagesResponseItem `json:"messages"`
}

// GetAllMessagesResponseItem is used in tandem with GetAllMessagesResponseData.
// It represents a single message. It has three fields: "id", "text", and
// "is_palindrome".
type GetAllMessagesResponseItem struct {
	ID           int    `json:"id"`
	Text         string `json:"text"`
	IsPalindrome *bool  `json:"is_palindrome"` // trinary, nil if unknown
}
