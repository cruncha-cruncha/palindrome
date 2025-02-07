package main

import (
	"testing"
	"time"
)

func TestMessageOrchestratorAdd(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, err, _, _ := mo.Add("hello")
	if err != nil {
		t.Fatalf(`mo.Add("hello") = %v, want nil`, err)
	}

	if id != 1 {
		t.Fatalf(`mo.Add("hello") = %d, want 1`, id)
	}
}

func TestMessageOrchestratorAddMultiple(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, err, _, _ := mo.Add("hello")
	if err != nil {
		t.Fatalf(`mo.Add("hello") = %v, want nil`, err)
	}

	if id != 1 {
		t.Fatalf(`mo.Add("hello") = %d, want 1`, id)
	}

	id, err, _, _ = mo.Add("goodbye")
	if err != nil {
		t.Fatalf(`mo.Add("goodbye") = %v, want nil`, err)
	}

	if id != 2 {
		t.Fatalf(`mo.Add("goodbye") = %d, want 2`, id)
	}
}

func TestMessageOrchestratorGet(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, _, _, _ := mo.Add("hello")
	message, err := mo.Get(id)
	if err != nil {
		t.Fatalf(`mo.Get(%d) = %v, want nil`, id, err)
	}

	if message.text != "hello" {
		t.Fatalf(`mo.Get(%d).Text = %s, want "hello"`, id, message.text)
	}
}

func TestMessageOrchestratorGetNone(t *testing.T) {
	mo := NewMessageOrchestrator()

	_, err := mo.Get(1)
	if err == nil {
		t.Fatalf(`mo.Get(1) = nil, want non-nil`)
	}
}

func TestMessageOrchestratorUpdate(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, _, _, _ := mo.Add("hello")
	err, _, _ := mo.Update(id, "goodbye")
	if err != nil {
		t.Fatalf(`mo.Update(%d, "goodbye") = %v, want nil`, id, err)
	}

	message, err := mo.Get(id)
	if err != nil {
		t.Fatalf(`mo.Get(%d) = %v, want nil`, id, err)
	}

	if message.text != "goodbye" {
		t.Fatalf(`mo.Get(%d).Text = %s, want "goodbye"`, id, message.text)
	}
}

func TestMessageOrchestratorUpdateNone(t *testing.T) {
	mo := NewMessageOrchestrator()

	err, _, _ := mo.Update(1, "goodbye")
	if err == nil {
		t.Fatalf(`mo.Update(1, "goodbye") = nil, want non-nil`)
	}
}

func TestMessageOrchestratorDelete(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, _, _, _ := mo.Add("hello")
	err := mo.Delete(id)
	if err != nil {
		t.Fatalf(`mo.Delete(%d) = %v, want nil`, id, err)
	}

	_, err = mo.Get(id)
	if err == nil {
		t.Fatalf(`mo.Get(%d) = nil, want non-nil`, id)
	}
}

func TestMessageOrchestratorDeleteNone(t *testing.T) {
	mo := NewMessageOrchestrator()

	err := mo.Delete(1)
	if err == nil {
		t.Fatalf(`mo.Delete(1) = nil, want non-nil`)
	}
}

func TestMessageOrchestratorGetAll(t *testing.T) {
	mo := NewMessageOrchestrator()

	id1, _, _, _ := mo.Add("hello")
	id2, _, _, _ := mo.Add("goodbye")
	messages, err := mo.GetAll()
	if err != nil {
		t.Fatalf(`mo.GetAll() = %v, want nil`, err)
	}

	if len(messages) != 2 {
		t.Fatalf(`len(mo.GetAll()) = %d, want 2`, len(messages))
	}

	if messages[0].id != id1 {
		t.Fatalf(`mo.GetAll()[0].id = %d, want %d`, messages[0].id, id1)
	}

	if messages[0].text != "hello" {
		t.Fatalf(`mo.GetAll()[0].text = %s, want "hello"`, messages[0].text)
	}

	if messages[1].id != id2 {
		t.Fatalf(`mo.GetAll()[1].id = %d, want %d`, messages[1].id, id2)
	}

	if messages[1].text != "goodbye" {
		t.Fatalf(`mo.GetAll()[1].text = %s, want "goodbye"`, messages[1].text)
	}
}

func TestMessageOrchestratorGetAllNone(t *testing.T) {
	mo := NewMessageOrchestrator()

	messages, err := mo.GetAll()
	if err != nil {
		t.Fatalf(`mo.GetAll() = %v, want nil`, err)
	}

	if len(messages) != 0 {
		t.Fatalf(`len(mo.GetAll()) = %d, want 0`, len(messages))
	}
}

func TestMessageOrchestratorCombined(t *testing.T) {
	mo := NewMessageOrchestrator()

	id1, _, _, _ := mo.Add("hello")
	id2, _, _, _ := mo.Add("goodbye")
	err := mo.Delete(id1)
	messages, err := mo.GetAll()

	if err != nil {
		t.Fatalf(`mo.GetAll() = %v, want nil`, err)
	}

	if len(messages) != 1 {
		t.Fatalf(`len(mo.GetAll()) = %d, want 1`, len(messages))
	}

	if messages[0].id != id2 {
		t.Fatalf(`mo.GetAll()[0].id = %d, want %d`, messages[0].id, id2)
	}

	if messages[0].text != "goodbye" {
		t.Fatalf(`mo.GetAll()[0].text = %s, want "goodbye"`, messages[0].text)
	}
}

func TestMessageOrchestrationIsPalindrome(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, _, _, _ := mo.Add("racecar")

	time.Sleep(50 * time.Millisecond)
	message, err := mo.Get(id)

	if err != nil {
		t.Fatalf(`mo.Get(%d) = %v, want nil`, id, err)
	}

	if message.isPalindrome != P_TRUE {
		t.Fatalf(`mo.Get(%d).isPalindrome = %d, want %d`, id, message.isPalindrome, P_TRUE)
	}
}

func TestMessageOrchestrationIsNotPalindrome(t *testing.T) {
	mo := NewMessageOrchestrator()

	id, _, _, _ := mo.Add("hello")

	time.Sleep(50 * time.Millisecond)
	message, err := mo.Get(id)

	if err != nil {
		t.Fatalf(`mo.Get(%d) = %v, want nil`, id, err)
	}

	if message.isPalindrome != P_FALSE {
		t.Fatalf(`mo.Get(%d).isPalindrome = %d, want %d`, id, message.isPalindrome, P_FALSE)
	}
}
