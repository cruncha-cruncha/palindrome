package main

import (
	"errors"
	"sync"
	"sync/atomic"
)



type MessageOrchestrator struct {
	messages sync.Map
	nextId   atomic.Uint64
}

func NewMessageOrchestrator() MessageOrchestrator {
	return MessageOrchestrator{
		messages: sync.Map{},
		nextId:   atomic.Uint64{},
	}
}

type Message struct {
	id           int // this is redundant, but also small and convenient
	text         string
	isPalindrome int
	stopChan     chan bool
}

// if palindrome calculation was truly slow, I would add a cache in here based on text hash, so duplicate text could be detected quickly

// returns id, stopChan (can send a message on this channel to cancel the palindrome calculation), and doneChan (will receive a message / close when the isPalindrome calculation is complete)
func (mo *MessageOrchestrator) Add(text string) (int, error, chan bool, chan bool) {
	// atomically read and increment nextId
	id := int(mo.nextId.Add(1))

	// create a stopChan
	stopChan := make(chan bool, 1)

	// make a new Message, with a stopChan
	msg := Message{
		id:           id,
		text:         text,
		isPalindrome: P_UNKNOWN,
		stopChan:     stopChan,
	}

	// save to mo.messages (locked)
	mo.messages.Store(id, msg)

	// create a doneChan
	doneChan := make(chan bool, 1)

	go func(msg Message, doneChan chan bool) {
		// calculate isPalindrome
		isPalindrome := StringIsPalindrome(text)

		// if this was a real slow calculation, I would break it up and check stopChan periodically

		select {
		case <-msg.stopChan:
			// if any value, then don't update
		default:
			// update only if no changes since we started
			mo.messages.CompareAndSwap(msg.id, msg, Message{
				id:           msg.id,
				text:         msg.text,
				isPalindrome: isPalindrome,
				stopChan:     nil,
			})
		}

		doneChan <- true
	}(msg, doneChan)

	return id, nil, stopChan, doneChan
}

func (mo *MessageOrchestrator) Get(id int) (Message, error) {
	if msg, ok := mo.messages.Load(id); !ok {
		return Message{}, errors.New("Message not found")
	} else {
		return msg.(Message), nil
	}
}

// returns success (existing item), stopChan (can send a message on this channel to cancel the palindrome calculation), and doneChan (will receive a message / close when the isPalindrome calculation is complete)
func (mo *MessageOrchestrator) Update(id int, text string) (error, chan bool, chan bool) {
	// create a new stopChan
	stopChan := make(chan bool, 1)

	msg := Message{
		id:           id,
		text:         text,
		isPalindrome: P_UNKNOWN,
		stopChan:     stopChan,
	}

	v, ok := mo.messages.Swap(id, msg)
	if !ok {
		// check if item exists, if not, then return false
		mo.messages.Delete(id)
		return errors.New("Nothing to update"), nil, nil
	}
	existing := v.(Message)

	// close existing stopChan if present
	if existing.stopChan != nil {
		existing.stopChan <- true
	}

	// create a doneChan
	doneChan := make(chan bool, 1)

	go func(msg Message, doneChan chan bool) {
		// calculate isPalindrome
		isPalindrome := StringIsPalindrome(text)

		// if this was a real slow calculation, I would break it up and check stopChan periodically

		select {
		case <-msg.stopChan:
			// if any value, then don't update
		default:
			// update only if no changes since we started
			mo.messages.CompareAndSwap(msg.id, msg, Message{
				id:           msg.id,
				text:         msg.text,
				isPalindrome: isPalindrome,
				stopChan:     nil,
			})
		}

		doneChan <- true
	}(msg, doneChan)

	return nil, stopChan, doneChan
}

func (mo *MessageOrchestrator) Delete(id int) error {
	_, ok := mo.messages.LoadAndDelete(id)
	if !ok {
		return errors.New("Message not found")
	} else {
		return nil
	}
}

func (mo *MessageOrchestrator) GetAll() ([]Message, error) {
	out := []Message{}
	mo.messages.Range(func(key, value any) bool {
		out = append(out, value.(Message))
		return true
	})

	return out, nil
}
