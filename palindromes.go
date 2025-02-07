package main

import (
	"errors"
	"sync"
)

type Palindromes struct {
	lock sync.RWMutex
	work map[string]PalindromeWork
}

func NewPalindromes() Palindromes {
	return Palindromes{
		lock: sync.RWMutex{},
		work: make(map[string]PalindromeWork),
	}
}

type PalindromeWork struct {
	hash     string
	status   PalindromeWorkStatus
	messages map[int]chan PalindromeWorkStatus
	cancel   chan bool
}

func (p *Palindromes) Add(msg Message) (key PalindromeWorkKey, current PalindromeWorkStatus, onChange chan PalindromeWorkStatus, err error) {
	key = PalindromeWorkKey{
		hash:      msg.hash,
		messageId: msg.id,
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	work, ok := p.work[msg.hash]

	if ok {
		onChange = make(chan PalindromeWorkStatus, 1)
		if listener, ok := work.messages[msg.id]; ok {
			onChange = listener
		} else {
			work.messages[msg.id] = onChange
		}

		return key, work.status, onChange, nil
	}

	work = PalindromeWork{
		hash:     msg.hash,
		messages: map[int]chan PalindromeWorkStatus{msg.id: make(chan PalindromeWorkStatus, 1)},
		status: PalindromeWorkStatus{
			isPalindrome: P_UNKNOWN,
			done:         false,
		},
		cancel: make(chan bool, 1),
	}
	p.work[msg.hash] = work

	go p.doWork(msg)

	return key, work.status, work.messages[0], nil
}

func (p *Palindromes) Remove(key PalindromeWorkKey) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if work, ok := p.work[key.hash]; ok {
		listener, ok := work.messages[key.messageId]
		if !ok {
			return errors.New("Not found")
		}

		close(listener)
		delete(work.messages, key.messageId)

		if len(work.messages) == 0 {
			delete(p.work, key.hash)
			select {
			case work.cancel <- true:
			default:
			}
			close(work.cancel)
		}
	}

	return nil
}

func (p *Palindromes) Poll(key PalindromeWorkKey) (found bool, current PalindromeWorkStatus, onChange chan PalindromeWorkStatus, err error) {
	p.lock.RLock()
	work, ok := p.work[key.hash]
	p.lock.RUnlock()

	if !ok {
		return false, PalindromeWorkStatus{}, nil, nil
	}

	if onChange, ok = work.messages[key.messageId]; !ok {
		return true, work.status, nil, nil
	} else {
		return true, work.status, onChange, nil
	}
}

func (p *Palindromes) Clear() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, work := range p.work {
		for _, listener := range work.messages {
			close(listener)
		}

		select {
		case work.cancel <- true:
		default:
		}
		close(work.cancel)
	}

	p.work = make(map[string]PalindromeWork)

	return nil
}
