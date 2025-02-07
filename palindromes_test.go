package main

import "testing"

func newFakeMessage() Message {
	return Message{
		id:   1,
		hash: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		text: "hello",
	}
}

func TestPalindromeOrchestratorAdd(t *testing.T) {
	po := NewPalindromes()

	msg := newFakeMessage()

	key, _, _, err := po.Add(msg)
	if err != nil {
		t.Fatalf(`po.Add(%+v) has err %+v, want nil`, msg, err)
	}
	if key.hash != msg.hash {
		t.Fatalf(`po.Add(%+v) key.hash = %s, want %s`, msg, key.hash, msg.hash)
	}
	if key.messageId != msg.id {
		t.Fatalf(`po.Add(%+v) key.messageId = %d, want %d`, msg, key.messageId, msg.id)
	}
}

func TestPalindromeOrchestratorRemove(t *testing.T) {
	po := NewPalindromes()

	msg := newFakeMessage()

	key, _, _, _ := po.Add(msg)

	err := po.Remove(key)
	if err != nil {
		t.Fatalf(`po.Remove(%+v) has err %+v, want nil`, key, err)
	}

	found, _, _, err := po.Poll(key)
	if err != nil {
		t.Fatalf(`po.Poll(%+v) has err %+v, want nil`, key, err)
	}
	if found {
		t.Fatalf(`po.Poll(%+v) found, want not found`, key)
	}
}

func TestPalindromeOrchestratorPoll(t *testing.T) {
	po := NewPalindromes()

	msg := newFakeMessage()

	key, _, _, _ := po.Add(msg)

	found, _, _, err := po.Poll(key)
	if err != nil {
		t.Fatalf(`po.Poll(%+v) has err %+v, want nil`, key, err)
	}
	if !found {
		t.Fatalf(`po.Poll(%+v) not found`, key)
	}
}

func TestPalindromeOrchestratorClear(t *testing.T) {
	po := NewPalindromes()

	msg := newFakeMessage()

	key, _, _, _ := po.Add(msg)

	err := po.Clear()
	if err != nil {
		t.Fatalf(`po.Clear() has err %+v, want nil`, err)
	}

	found, _, _, err := po.Poll(key)
	if err != nil {
		t.Fatalf(`po.Poll(%+v) has err %+v, want nil`, key, err)
	}
	if found {
		t.Fatalf(`po.Poll(%+v) found, want not found`, key)
	}
}
