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
	Text string `json:"text"`
	IsPalindrome *bool `json:"is_palindrome"`
}

type GetAllMessagesResponseData struct {
	Messages []GetAllMessagesResponseItem `json:"messages"`
}

type GetAllMessagesResponseItem struct {
	ID int `json:"id"`
	Text string `json:"text"`
	IsPalindrome bool `json:"is_palindrome"`
}

// conveniece functions for converting between types

// prefer this over a struct method
func NewGetMessageResponseDataFromMessage(m *Message) GetMessageResponseData {
	return GetMessageResponseData{
		Text: m.text,
		IsPalindrome: PalindromeStatusToBool(m.isPalindrome),
	}
}

func NewGetMessagesResponseItemFromMessage(m *Message) GetAllMessagesResponseItem {
	return GetAllMessagesResponseItem{
		ID: m.id,
		Text: m.text,
		IsPalindrome: *PalindromeStatusToBool(m.isPalindrome),
	}
}