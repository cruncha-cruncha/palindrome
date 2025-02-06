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
	stopChan chan bool
}

// returns id, stopChan, and a channel that will be closed when the isPalindrome calculation is complete
func (mo *MessageOrchestrator) Add(text string) (ID, chan bool, chan bool) {
	// atomically read and increment ss.nextId
	// create hash of text
	// make a new Message, with a stopChan
	// save to ss.messages

	// create a doneChan
	// spawn a goroutine to: calculate isPalindrome (check stopChan periodically), update (locked) or don't update if hash has changed / item deleted, then close the doneChan

	return 0, nil, nil
}

func (mo *MessageOrchestrator) Get(id ID) (Message, bool) {
	// read (locked) from ss.messages
	// return the Message and a bool indicating success

	return Message{}, false
}

// returns success, stopChan, and a channel that will be closed when the isPalindrome calculation is complete
func (mo *MessageOrchestrator) Update(id ID, text string) (bool, chan bool, chan bool) {
	// calculate the new text hash
	// create a new stopChan
	// obtain a write lock
	// check if item exists, if not, then return false
	// close existing stopChan if present
	// update text, hash, isPalindrome, stopChan, then drop the lock
	// read (locked) from ss.messages

	// create a doneChan
	// spawn a goroutine to: calculate isPalindrome (check stopChan periodically), update (locked) or don't update if hash has changed / item deleted, then close the doneChan

	return false, nil, nil
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