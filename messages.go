package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Messages implements MessageOrchestrator. It stores messages in-memory (is not
// persistent). It is safe for concurrent use. The first message added will be
// assigned an id of 1, then 2, then 3, etc. Ids are not reused.
type Messages struct {
	messages sync.Map
	nextId   atomic.Uint64
}

// NewMessageOrchestrator creates a new MessageOrchestrator with no messages.
func NewMessages() Messages {
	return Messages{
		messages: sync.Map{},
		nextId:   atomic.Uint64{},
	}
}

func (m *Messages) Add(text string) (Message, error) {
	msg := Message{
		id:   int(m.nextId.Add(1)),
		hash: CalculateHash(text),
		text: text,
	}

	m.messages.Store(msg.id, msg)

	// go func(msg Message, done chan int) {
	// 	// calculate isPalindrome
	// 	isPalindrome := StringIsPalindrome(text)

	// 	// if this was a real slow calculation, I would break it up and check stopChan periodically

	// 	select {
	// 	case <-msg.stop:
	// 		// if any value, then don't update
	// 	default:
	// 		// update only if no changes since we started
	// 		mo.messages.CompareAndSwap(msg.id, msg, Message{
	// 			id:           msg.id,
	// 			text:         msg.text,
	// 			isPalindrome: isPalindrome,
	// 			stop:         nil,
	// 		})
	// 	}

	// 	done <- isPalindrome
	// }(msg, done)

	return msg, nil
}

func (m *Messages) Get(id int) (Message, bool, error) {
	if msg, ok := m.messages.Load(id); !ok {
		return Message{}, false, nil
	} else {
		return msg.(Message), true, nil
	}
}

func (m *Messages) Update(id int, text string) (Message, error) {
	msg := Message{
		id:   id,
		hash: CalculateHash(text),
		text: text,
	}

	_, ok := m.messages.Swap(id, msg)
	if !ok {
		// check if item exists, if not, then return false
		m.messages.Delete(id)
		return Message{}, errors.New("Nothing to update")
	}

	return msg, nil
}

func (m *Messages) Delete(id int) error {
	m.messages.LoadAndDelete(id)
	return nil
}

func (m *Messages) GetAll() ([]Message, error) {
	out := []Message{}
	m.messages.Range(func(key, value any) bool {
		out = append(out, value.(Message))
		return true
	})

	return out, nil
}

func (m *Messages) DeleteAll() error {
	m.messages.Clear()
	return nil
}
