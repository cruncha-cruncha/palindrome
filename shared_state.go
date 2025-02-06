package main

// I like this shared "server struct" pattern better than wrapping all handlers in closures
type SharedState struct {
	mo MessageOrchestration
}

func NewSharedState() SharedState {
	mo := NewMessageOrchestrator()
	return SharedState{
		mo: &mo,
	}
}
