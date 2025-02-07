package main

// import (
// 	"testing"
// )

// func TestStringIsPalindromeCaseOne(t *testing.T) {
// 	if StringIsPalindrome("racecar") != P_TRUE {
// 		t.Fatalf(`StringIsPalindrome("racecar") = %d, want %d`, StringIsPalindrome("racecar"), P_TRUE)
// 	}
// }

// func TestStringIsPalindromeCaseTwo(t *testing.T) {
// 	if StringIsPalindrome("") != P_UNKNOWN {
// 		t.Fatalf(`StringIsPalindrome("") = %d, want %d`, StringIsPalindrome(""), P_UNKNOWN)
// 	}
// }

// func TestStringIsPalindromeCaseThree(t *testing.T) {
// 	if StringIsPalindrome("hello") != P_FALSE {
// 		t.Fatalf(`StringIsPalindrome("hello") = %d, want %d`, StringIsPalindrome("hello"), P_FALSE)
// 	}
// }

// func TestStringIsPalindromeCaseFour(t *testing.T) {
// 	if StringIsPalindrome("A man, a plan, a canal, panama!") != P_TRUE {
// 		t.Fatalf(`StringIsPalindrome("A man, a plan, a canal, panama!") = %d, want %d`, StringIsPalindrome("A man, a plan, a canal, panama!"), P_TRUE)
// 	}
// }

// func TestPalindromeStatusToBoolPointerCaseOne(t *testing.T) {
// 	result := PalindromeStatusToBoolPointer(P_TRUE)
// 	if result == nil {
// 		t.Fatalf(`PalindromeStatusToBoolPointer(P_TRUE) = nil, want non-nil`)
// 	} else if *result != true {
// 		t.Fatalf(`*PalindromeStatusToBoolPointer(P_TRUE) = %t, want %t`, *PalindromeStatusToBoolPointer(P_TRUE), true)
// 	}
// }

// func TestPalindromeStatusToBoolPointerCaseTwo(t *testing.T) {
// 	if PalindromeStatusToBoolPointer(P_UNKNOWN) != nil {
// 		t.Fatalf(`PalindromeStatusToBoolPointer(P_UNKNOWN) = %v, want nil`, PalindromeStatusToBoolPointer(P_UNKNOWN))
// 	}
// }

// func TestPalindromeStatusToBoolPointerCaseThree(t *testing.T) {
// 	result := PalindromeStatusToBoolPointer(P_FALSE)
// 	if result == nil {
// 		t.Fatalf(`PalindromeStatusToBoolPointer(P_FALSE) = nil, want non-nil`)
// 	} else if *result != false {
// 		t.Fatalf(`*PalindromeStatusToBoolPointer(P_FALSE) = %t, want %t`, *PalindromeStatusToBoolPointer(P_FALSE), false)
// 	}
// }
