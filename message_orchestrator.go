package main

type ID int

type MessageOrchestrator struct {
	messages map[int]Message
	nextId ID
}

// could use composition here, but I'd rather not
type Message struct {
	id int // this is redundant, but also small and convenient
	hash string // a hash of the text, for quick comparison
	text string
	isPalindrome int
}

func (mo *MessageOrchestrator) Add(text string) (ID, chan bool) {
	// atomically read and increment ss.nextId
	// create hash of text
	// make a new Message
	// save to ss.messages

	// create a channel
	// spawn a goroutine to: calculate isPalindrome, update (locked) or don't update if hash has changed / item deleted, then close the channel

	return 0, nil
}

func (mo *MessageOrchestrator) Get(id ID) (Message, bool) {
	// read (locked) from ss.messages
	// return the Message and a bool indicating success

	return Message{}, false
}

func (mo *MessageOrchestrator) Update(id ID, text string) (bool, chan bool) {
	// calculate the new text hash
	// obtain a write lock
	// check if item exists, if not, then return false
	// update text and hash, then drop the lock
	// read (locked) from ss.messages

	// create a channel
	// spawn a goroutine to: calculate isPalindrome, update (locked) or don't update if hash has changed / item deleted, then close the channel

	return false, nil
}

func (mo *MessageOrchestrator) Delete(id ID) bool {
	// delete (locked) from ss.messages

	return false
}

func (mo *MessageOrchestrator) GetAll() []Message {
	// read (locked) from ss.messages
	// return a slice of all messages

	return nil
}