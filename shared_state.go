package main

// SharedState contains all the information that a handler might need: every
// handler is a method on this struct. As such, all fields and operations must
// be safe for concurrent use.
//
// I like this shared "server struct" pattern better than wrapping all handlers
// in closures, because I find having all shared state in one place makes it
// easier to understand a service at a glance and reduces boilerplate code.
type SharedState struct {
	mo MessageOrchestrator
	po WorkOrchestrator[Message, PalindromeWorkKey, PalindromeWorkStatus]
}

// MessageOrchestration is an interface for a service that can store and
// manipulate messages. It's a simple abstraction that allows us to swap out
// the underlying implementation (maybe switching to a database) without
// changing the rest of the code.
type MessageOrchestrator interface {
	Add(text string) (Message, error)
	Get(id int) (Message, bool, error)
	Update(id int, text string) (Message, error)
	Delete(id int) error
	GetAll() ([]Message, error)
	DeleteAll() error
}

type Message struct {
	id   int
	hash string
	text string
}

type WorkOrchestrator[D any, K any, R any] interface {
	Add(d D) (key K, current R, onChange chan R, err error)
	Remove(key K) error
	Poll(key K) (found bool, current R, onChange chan R, err error)
	Clear() error
}

type PalindromeWorkStatus struct {
	isPalindrome int
	done         bool
}

type PalindromeWorkKey struct {
	hash      string
	messageId int
}

// NewSharedState initializes all fields so they're ready to use. It should be
// called once at the beginning of the program.
func NewSharedState() SharedState {
	mo := NewMessages()
	po := NewPalindromes()

	return SharedState{
		mo: &mo,
		po: &po,
	}
}

func PWorkKeyFromMsg(msg Message) PalindromeWorkKey {
	return PalindromeWorkKey{
		hash:      msg.hash,
		messageId: msg.id,
	}
}
