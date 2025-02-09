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

// NewMessages creates a new Messages struct with no messages.
func NewMessages() Messages {
	return Messages{
		messages: sync.Map{},
		nextId:   atomic.Uint64{},
	}
}

// Add takes in some text and returns a Message, with a unique id and the hash
// of that text. This particular implementation will never throw an error. Once
// a message is added, it's immediately available for retrieval / deletion.
func (m *Messages) Add(text string) (Message, error) {
	msg := Message{
		id:   int(m.nextId.Add(1)),
		hash: CalculateHash(text),
		text: text,
	}

	m.messages.Store(msg.id, msg)

	return msg, nil
}

// Get returns a Message by id. This particular implementation will never throw
// an error, but it will return false if the message doesn't exist.
func (m *Messages) Get(id int) (Message, bool, error) {
	if msg, ok := m.messages.Load(id); !ok {
		return Message{}, false, nil
	} else {
		return msg.(Message), true, nil
	}
}

// Update takes in a Message id and some text. It will completely replace the 
// corresponding Message's text and update it's hash if the Message exists. If 
// not, it will throw and error.
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

// Delete removes a Message by id. This particular implementation will never
// throw an error. There is no way to tell if the message existed or not.
func (m *Messages) Delete(id int) error {
	m.messages.LoadAndDelete(id)
	return nil
}

// GetAll returns all messages in the system. This particular implementation
// will never throw an error. Due to the limitations of the sync.Map type, the
// messages are not guaranteed to be in any particular order nor do they
// represent a single snapshot at one point in time.
func (m *Messages) GetAll() ([]Message, error) {
	out := []Message{}
	m.messages.Range(func(key, value any) bool {
		out = append(out, value.(Message))
		return true
	})

	return out, nil
}

// DeleteAll removes all messages from the system. This particular
// implementation will never throw an error.
func (m *Messages) DeleteAll() error {
	m.messages.Clear()
	return nil
}
