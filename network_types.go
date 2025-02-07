package main

// request types

type CreateMessageRequestData struct {
	Text string `json:"text"`
}

type UpdateMessageRequestData struct {
	Text string `json:"text"`
}

// response types

type CreateMessageResponseData struct {
	ID int `json:"id"`
}

type GetMessageResponseData struct {
	Text         string `json:"text"`
	IsPalindrome *bool  `json:"is_palindrome"` // trinary, nil if unknown
	// in actual production code, I would use an explicit status field instead, this is just for fun
}

type GetAllMessagesResponseData struct {
	Messages []GetAllMessagesResponseItem `json:"messages"`
}

type GetAllMessagesResponseItem struct {
	ID           int    `json:"id"`
	Text         string `json:"text"`
	IsPalindrome *bool  `json:"is_palindrome"` // trinary, nil if unknown
}

// convenience functions for converting between types

func NewMsgResponseData(m *Message) GetMessageResponseData {
	f := false
	return GetMessageResponseData{
		Text:         m.text,
		IsPalindrome: &f,
	}
}

func NewMsgsResponseItem(m *Message) GetAllMessagesResponseItem {
	f := false
	return GetAllMessagesResponseItem{
		ID:           m.id,
		Text:         m.text,
		IsPalindrome: &f,
	}
}
