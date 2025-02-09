package main

import (
	"errors"
	"sync"
)

// Palindromes implements WorkOrchestrator. The "work" it does is determining if
// a string is a palindrome. It stores everything in-memory (is not persistent).
// It's safe for concurrent use.
// 
// If two messages have the same text, they will have the same hash, and share
// the same PalindromeWork. They will each have their own listener (a channel
// which receives a message everytime result changes). If all messages with the
// same hash are removed, the corresponding PalindromeWork is removed. Old work
// is not cached.
type Palindromes struct {
	lock sync.RWMutex
	work map[string]PalindromeWork
}

// NewPalindromes creates a new Palindromes struct with no work.
func NewPalindromes() Palindromes {
	return Palindromes{
		lock: sync.RWMutex{},
		work: make(map[string]PalindromeWork),
	}
}

// PalindromeWork holds all the information necessary to determine if a string
// is a palindrome. The result field contains the result of the calculation 
// and whether or not the calculation is done. Listeners and cancel should never
// be closed outside of Palindromes methods.
type PalindromeWork struct {
	hash      string
	result    PWResult
	// key: message id, value: a channel, receives updates when result changes
	listeners map[int]chan PWResult
	// Used to abort work early
	cancel    chan bool
}

// Add takes in a Message, kicks off work on calculating if it's a palindrome 
// (if work hasn't already started / been completed), and returns a
// PalindromeWorkKey (which can be used to delete work or poll progress), the
// current state of work (may be complete or in progress), a channel which will
// receive updates when the state of work changes, and an error. In practice,
// this method will never error. The onChange channel is unique per message id.
// This method is safe for concurrent use. If there is work to do, it calls
// doWork in new a goroutine.
func (p *Palindromes) Add(msg Message) (key PWKey, current PWResult, onChange chan PWResult, err error) {
	key = PWKey{
		hash:      msg.hash,
		messageId: msg.id,
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	work, ok := p.work[msg.hash]

	if ok {
		onChange = make(chan PWResult, 1)
		if listener, ok := work.listeners[msg.id]; ok {
			onChange = listener
		} else {
			work.listeners[msg.id] = onChange
		}

		return key, work.result, onChange, nil
	}

	work = PalindromeWork{
		hash:     msg.hash,
		listeners: map[int]chan PWResult{msg.id: make(chan PWResult, 1)},
		result: PWResult{
			isPalindrome: P_UNKNOWN,
			done:         false,
		},
		cancel: make(chan bool, 1),
	}
	p.work[msg.hash] = work

	go p.doWork(msg)

	return key, work.result, work.listeners[0], nil
}

// Remove is used to cancel or delete work. If work is in progress and no other
// messages are relying on it (aka no listeners), then the work is cancelled and
// removed. If no work is found, no action is taken. If work is found but other
// messages are relying on it, only this message's listener is removed. This
// method is safe for concurrent use.
func (p *Palindromes) Remove(key PWKey) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if work, ok := p.work[key.hash]; ok {
		listener, ok := work.listeners[key.messageId]
		if !ok {
			return errors.New("Not found")
		}

		close(listener)
		delete(work.listeners, key.messageId)

		if len(work.listeners) == 0 {
			delete(p.work, key.hash)
			// write asynchronously
			select {
			case work.cancel <- true:
			default:
			}
			close(work.cancel)
		}
	}

	return nil
}

// Poll is used to check on the progress of work, and possibly get the resulting
// value. Even if work is not done, can listen to onChange for updates. This
// method is safe for concurrent use.
//
// If work corresponding to the key's hash is found, but there is no listener
// for the key's messageId, then found is true but onChange is nil. No listener
// is added.
func (p *Palindromes) Poll(key PWKey) (found bool, current PWResult, onChange chan PWResult, err error) {
	p.lock.RLock()
	work, ok := p.work[key.hash]
	p.lock.RUnlock()

	if !ok {
		return false, PWResult{}, nil, nil
	}

	if onChange, ok = work.listeners[key.messageId]; !ok {
		return true, work.result, nil, nil
	} else {
		return true, work.result, onChange, nil
	}
}

// Clear is used to immediately cancel and remove all work and listeners. This
// method is safe for concurrent use.
func (p *Palindromes) Clear() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, work := range p.work {
		for _, listener := range work.listeners {
			close(listener)
		}

		// write asynchronously
		select {
		case work.cancel <- true:
		default:
		}
		close(work.cancel)
	}

	p.work = make(map[string]PalindromeWork)

	return nil
}
