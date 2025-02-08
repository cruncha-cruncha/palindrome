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

// Message is a simple struct for storing a message. It has three fields: an id
// (integer, unique, ascending), a hash (string, calculated from the text,
// hopefully unique), and the text (string, provided by the user). On adding a
// message to Messages, all three fields will be populated.
//
// Hash is used to de-duplicate work when calculating palindromes. If two
// messages have the same text, then they will have the same hash, and so only
// one palindrome calculation needs to be done.
type Message struct {
	id   int
	hash string
	text string
}

// WorkOrchestrator is an interface for helping manage long-running tasks, all
// of the same type (like calculating if a string is a palindrome). It's types
// are; D: all the Data needed to start work; K: a Key to identify any one
// piece of work; and R: the Result of some work. R should provide some
// indication of started/in progress/done. Once work is finished, the result is
// stored until explicitly removed. 
type WorkOrchestrator[D any, K any, R any] interface {
	// Add takes in some data and starts work on it. It returns a key which can
	// be used to cancel the work and remove it's result, or poll for progress.
	// Current result after just starting work is usually empty. OnChange will
	// recieve updates when the result changes.
	Add(d D) (key K, current R, onChange chan R, err error)
	// Remove cancels work and removes the result.
	Remove(key K) error
	// Poll returns the current result of work (could be in progress or done).
	// It also returns onChange which will recieve updates when the result
	// changes.
	Poll(key K) (found bool, current R, onChange chan R, err error)
	// Clear cancels all work and removes all results.
	Clear() error
}

// PalindromeWorkStatus is the result of a palindrome calculation. It has two
// fields: isPalindrome (P_UNKNOWN, P_TRUE, or P_FALSE) and done (bool).
type PalindromeWorkStatus struct {
	isPalindrome int
	done         bool
}

// PalindromeWorkKey is the unique identifier for a piece of palindrome
// calculation work. It has two fields: hash (string, hopefully unique to some
// text) and messageId (integer, unique to a message). Hash determines
// isPalindrome, while messageId determines onChange (each message gets its
// own listener).
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
