package main

import (
	"testing"
)

func TestEmptyStringIsPalindrome(t *testing.T) {
	text := ""	
	result := StringIsPalindrome(text)
	if result != P_UNKNOWN {
		t.Fatalf(`StringIsPalindrome(%v) = %d, want %d`, text, result, P_UNKNOWN)
	}
}

// racecar
func TestStringIsPalindromeCaseOne(t *testing.T) {
	text := "racecar"
	result := StringIsPalindrome(text)
	if result != P_TRUE {
		t.Fatalf(`StringIsPalindrome(%v) = %d, want %d`, text, result, P_TRUE)
	}
}

// hello
func TestStringIsPalindromeCaseTwo(t *testing.T) {
	text := "hello"
	result := StringIsPalindrome(text)
	if result != P_FALSE {
		t.Fatalf(`StringIsPalindrome(%v) = %d, want %d`, text, result, P_FALSE)
	}
}

// A man, a plan, a canal, panama!
func TestStringIsPalindromeCaseThree(t *testing.T) {
	text := "A man, a plan, a canal, panama!"
	result := StringIsPalindrome(text)
	if result != P_TRUE {
		t.Fatalf(`StringIsPalindrome(%v) = %d, want %d`, text, result, P_TRUE)
	}
}

func TestPalindromeStatusToBoolPointerCaseTrue(t *testing.T) {
	result := PStatusToBoolPointer(P_TRUE)
	if result == nil {
		t.Fatalf(`PStatusToBoolPointer(P_TRUE) = nil, want non-nil`)
	} else if *result != true {
		t.Fatalf(`*PStatusToBoolPointer(P_TRUE) = %t, want %t`, *result, true)
	}
}

func TestPalindromeStatusToBoolPointerCaseUnknown(t *testing.T) {
	result := PStatusToBoolPointer(P_UNKNOWN)
	if result != nil {
		t.Fatalf(`PStatusToBoolPointer(P_UNKNOWN) = %v, want nil`, result)
	}
}

func TestPalindromeStatusToBoolPointerCaseFalse(t *testing.T) {
	result := PStatusToBoolPointer(P_FALSE)
	if result == nil {
		t.Fatalf(`PStatusToBoolPointer(P_FALSE) = nil, want non-nil`)
	} else if *result != false {
		t.Fatalf(`*PStatusToBoolPointer(P_FALSE) = %t, want %t`, *result, false)
	}
}
