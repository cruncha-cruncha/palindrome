package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// P_UNKNOWN is used for both an empty string, and while calculating if a string
// is a palindrome. P_TRUE and P_FALSE are self-explanatory.
const (
	P_UNKNOWN = 0
	P_TRUE    = 1
	P_FALSE   = 2
)

// PStatusToBoolPointer converts P_UNKNOWN to nil, P_TRUE to true, and P_FALSE
// to false.
func PStatusToBoolPointer(status int) *bool {
	var out *bool

	if status == P_UNKNOWN {
		out = nil
	} else if status == P_TRUE {
		t := true
		out = &t
	} else {
		f := false
		out = &f
	}

	return out
}

// StringIsPalindrome returns P_UNKNOWN if the string is empty, P_TRUE if it is
// a palindrome, and P_FALSE if it is not. It's case-insensitive and only
// considers alphanumeric characters (whitespace and punctuation are ignored).
func StringIsPalindrome(s string) int {
	if len(s) == 0 {
		return P_UNKNOWN
	}

	// normalize string
	s = strings.ToLower(s)
	reg, _ := regexp.Compile(`[^\w\n]+`)
	s = reg.ReplaceAllString(s, "")

	length := len(s)

	// this was entirely generated by CoPilot
	for i := 0; i < length/2; i++ {
		if s[i] != s[length-i-1] {
			return P_FALSE
		}
	}

	return P_TRUE
}

// doWork is a Palindromes method that calculates if a message is a palindrome.
// Once completed, it saves the result and updates all listeners. It's safe to
// to run concurrently.
// 
// doWork can be artificially slowed down, and will take as long as S_DELAY
// seconds (default 0) to complete. It's also cancellable, and checks 4 times
// during S_DELAY to see if it should stop early. If stopped early, it leaves
// Palindromes as-is and does not send any updates to listeners.
func (p *Palindromes) doWork(msg Message) {
	isPalindrome := StringIsPalindrome(msg.text)

	newResult := PWResult{
		isPalindrome: isPalindrome,
		done:         true,
	}

	// pretend this is really slow
	var delay float64 = 0
	delay_str := os.Getenv("S_DELAY")
	if v, err := strconv.Atoi(delay_str); err == nil && v > 0 {
		delay = float64(v)
	}

	if delay > 0 {
		// check if we should stop work early, four times during the articial delay
		ms_delay := delay / 4 * 1000
		for i := 0; i < 4; i++ {
			time.Sleep(time.Duration(ms_delay) * time.Millisecond)

			p.lock.Lock()
			work, ok := p.work[msg.hash]
			p.lock.Unlock()

			if !ok {
				return
			} else {
				// write asynchronously
				select {
				case <-work.cancel:
					return
				default:
				}
			}
		}
	}

	// update work
	p.lock.Lock()
	if work, ok := p.work[msg.hash]; ok {
		work.result = newResult
		p.work[msg.hash] = work
		for _, listener := range work.listeners {
			// write asynchronously
			select {
			case listener <- newResult:
			default:
			}
		}
	}
	p.lock.Unlock()
}
