package main

import (
	"testing"
)

func TestMessageOrchestratorAdd(t *testing.T) {
	mo := NewMessages()

	text := "hello"
	msg, err := mo.Add(text)
	if err != nil {
		t.Fatalf(`mo.Add(%v) has err %+v, want nil`, text, err)
	}

	if msg.id != 1 {
		t.Fatalf(`mo.Add(%v) msg.id = %d, want 1`, text, msg.id)
	}
}

func TestMessageOrchestratorAddMultiple(t *testing.T) {
	mo := NewMessages()

	text := "hello"
	msg, err := mo.Add(text)
	if err != nil {
		t.Fatalf(`mo.Add(%v) has err %+v, want nil`, text, err)
	}

	if msg.id != 1 {
		t.Fatalf(`mo.Add(%v) = %d, want 1`, text, msg.id)
	}

	text = "goodbye"
	msg, err = mo.Add(text)
	if err != nil {
		t.Fatalf(`mo.Add(%v) has err %+v, want nil`, text, err)
	}

	if msg.id != 2 {
		t.Fatalf(`mo.Add("goodbye") msg.id = %d, want 2`, msg.id)
	}
}

func TestMessageOrchestratorGet(t *testing.T) {
	mo := NewMessages()

	original, _ := mo.Add("hello")
	msg, found, err := mo.Get(original.id)
	if err != nil {
		t.Fatalf(`mo.Get(%d) has err %+v, want nil`, original.id, err)
	}

	if !found {
		t.Fatalf(`mo.Get(%d) not found`, original.id)
	}

	if msg.text != original.text {
		t.Fatalf(`mo.Get(%d) msg.text = %s, want %s`, original.id, msg.text, original.text)
	}
}

func TestMessageOrchestratorGetNone(t *testing.T) {
	mo := NewMessages()

	_, found, err := mo.Get(1)
	if err != nil {
		t.Fatalf(`mo.Get(1) has err %+v, want nil`, err)
	}
	if found {
		t.Fatalf(`mo.Get(1) found, want not found`)
	}
}

func TestMessageOrchestratorUpdate(t *testing.T) {
	mo := NewMessages()

	original, _ := mo.Add("hello")

	msg, err := mo.Update(original.id, "goodbye")
	if err != nil {
		t.Fatalf(`mo.Update(%d, %v) has err %+v, want nil`, original.id, msg.text, err)
	}

	if msg.id != original.id {
		t.Fatalf(`mo.Update(%d, %v) msg.id = %d, want %d`, original.id, msg.text, msg.id, original.id)
	}
	if msg.hash == original.hash {
		t.Fatalf(`mo.Update(%d, %v) msg.hash = %s, want not %s`, original.id, msg.text, msg.hash, original.hash)
	}

	another, found, err := mo.Get(msg.id)
	if err != nil {
		t.Fatalf(`mo.Get(%d) has err %+v, want nil`, msg.id, err)
	}
	if !found {
		t.Fatalf(`mo.Get(%d) not found`, msg.id)
	}
	if another.text != msg.text {
		t.Fatalf(`mo.Get(%d) msg.Text = %s, want %s`, msg.id, another.text, msg.text)
	}
}

func TestMessageOrchestratorUpdateNone(t *testing.T) {
	mo := NewMessages()

	text := "huh"
	_, err := mo.Update(1, text)
	if err == nil {
		t.Fatalf(`mo.Update(1, %s) has no err, it should`, text)
	}
}

func TestMessageOrchestratorDelete(t *testing.T) {
	mo := NewMessages()

	msg, _ := mo.Add("hello")
	err := mo.Delete(msg.id)
	if err != nil {
		t.Fatalf(`mo.Delete(%d) has err %+v, want nil`, msg.id, err)
	}

	_, found, err := mo.Get(msg.id)
	if err != nil {
		t.Fatalf(`mo.Get(%d) has err %+v, want nil`, msg.id, err)
	}
	if found {
		t.Fatalf(`mo.Get(%d) found, want not found`, msg.id)
	}
}

func TestMessageOrchestratorDeleteNone(t *testing.T) {
	mo := NewMessages()

	err := mo.Delete(1)
	if err != nil {
		t.Fatalf(`mo.Delete(1) has err %+v, want nil`, err)
	}
}

func TestMessageOrchestratorGetAll(t *testing.T) {
	mo := NewMessages()

	msg1, _ := mo.Add("hello")
	msg2, _ := mo.Add("goodbye")
	messages, err := mo.GetAll()
	if err != nil {
		t.Fatalf(`mo.GetAll() has err %+v, want nil`, err)
	}

	if len(messages) != 2 {
		t.Fatalf(`len(mo.GetAll()) = %d, want 2`, len(messages))
	}

	if messages[0].id != msg1.id {
		t.Fatalf(`mo.GetAll()[0].id = %d, want %d`, messages[0].id, msg1.id)
	}

	if messages[0].text != msg1.text {
		t.Fatalf(`mo.GetAll()[0].text = %s, want %s`, messages[0].text, msg1.text)
	}

	if messages[1].id != msg2.id {
		t.Fatalf(`mo.GetAll()[1].id = %d, want %d`, messages[1].id, msg2.id)
	}

	if messages[1].text != msg2.text {
		t.Fatalf(`mo.GetAll()[1].text = %s, want %s`, messages[1].text, msg2.text)
	}
}

func TestMessageOrchestratorGetAllNone(t *testing.T) {
	mo := NewMessages()

	messages, err := mo.GetAll()
	if err != nil {
		t.Fatalf(`mo.GetAll() = %v, want nil`, err)
	}

	if len(messages) != 0 {
		t.Fatalf(`len(mo.GetAll()) = %d, want 0`, len(messages))
	}
}

func TestMessageOrchestrationAddDuplicateText(t *testing.T) {
	mo := NewMessages()

	msg1, _ := mo.Add("hello")
	msg2, _ := mo.Add("hello")

	if msg1.id == msg2.id {
		t.Fatalf(`mo.Add("hello") = %d, want %d`, msg1.id, msg2.id)
	}
}