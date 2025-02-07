package main

// I like this shared "server struct" pattern better than wrapping all handlers in closures
type SharedState struct {
	mo MessageOrchestration
}

type MessageOrchestration interface {
	Add(text string) (int, error, chan bool, chan bool)
	Get(id int) (Message, error)
	Update(id int, text string) (error, chan bool, chan bool)
	Delete(id int) error
	GetAll() ([]Message, error)
	DeleteAll() error
}

func NewSharedState() SharedState {
	mo := NewMessageOrchestrator()
	return SharedState{
		mo: &mo,
	}
}
