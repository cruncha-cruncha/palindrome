package main

import (
	"slices"
	"testing"
)

func TestBinaryInsertionSortCaseOne(t *testing.T) {
	unsorted := []Message{
		{
			id:   1,
			text: "one",
		}, {
			id:   3,
			text: "three",
		}, {
			id:   2,
			text: "two",
		},
	}

	sorted := []Message{}

	for _, m := range unsorted {
		insertIndex := BinarySearch(sorted, func(m *Message) int { return m.id }, m.id)
		sorted = slices.Insert(sorted, insertIndex, Message{
			id:   m.id,
			text: m.text,
		})
	}

	if len(sorted) != len(unsorted) {
		t.Fatalf(`len(sorted) = %d, want %d`, len(sorted), len(unsorted))
	}

	if sorted[0].id != 1 {
		t.Fatalf(`sorted[0].id = %d, want 1`, sorted[0].id)
	}

	if sorted[1].id != 2 {
		t.Fatalf(`sorted[1].id = %d, want 2`, sorted[1].id)
	}

	if sorted[2].id != 3 {
		t.Fatalf(`sorted[2].id = %d, want 3`, sorted[2].id)
	}
}

func TestBinaryInsertionSortCaseTwo(t *testing.T) {
	unsorted := []Message{
		{
			id:   1,
			text: "one",
		}, {
			id:   2,
			text: "two",
		}, {
			id:   3,
			text: "three",
		},
	}

	sorted := []Message{}

	for _, m := range unsorted {
		insertIndex := BinarySearch(sorted, func(m *Message) int { return m.id }, m.id)
		sorted = slices.Insert(sorted, insertIndex, Message{
			id:   m.id,
			text: m.text,
		})
	}

	if len(sorted) != len(unsorted) {
		t.Fatalf(`len(sorted) = %d, want %d`, len(sorted), len(unsorted))
	}

	if sorted[0].id != 1 {
		t.Fatalf(`sorted[0].id = %d, want 1`, sorted[0].id)
	}

	if sorted[1].id != 2 {
		t.Fatalf(`sorted[1].id = %d, want 2`, sorted[1].id)
	}

	if sorted[2].id != 3 {
		t.Fatalf(`sorted[2].id = %d, want 3`, sorted[2].id)
	}
}

func TestBinaryInsertionSortCaseThree(t *testing.T) {
	unsorted := []Message{
		{
			id:   3,
			text: "three",
		}, {
			id:   2,
			text: "two",
		}, {
			id:   1,
			text: "one",
		},
	}

	sorted := []Message{}

	for _, m := range unsorted {
		insertIndex := BinarySearch(sorted, func(m *Message) int { return m.id }, m.id)
		sorted = slices.Insert(sorted, insertIndex, Message{
			id:   m.id,
			text: m.text,
		})
	}

	if len(sorted) != len(unsorted) {
		t.Fatalf(`len(sorted) = %d, want %d`, len(sorted), len(unsorted))
	}

	if sorted[0].id != 1 {
		t.Fatalf(`sorted[0].id = %d, want 1`, sorted[0].id)
	}

	if sorted[1].id != 2 {
		t.Fatalf(`sorted[1].id = %d, want 2`, sorted[1].id)
	}

	if sorted[2].id != 3 {
		t.Fatalf(`sorted[2].id = %d, want 3`, sorted[2].id)
	}
}
