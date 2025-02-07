package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	P_UNKNOWN = 0
	P_TRUE    = 1
	P_FALSE   = 2
)

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

func (p *Palindromes) doWork(msg Message) {
	isPalindrome := StringIsPalindrome(msg.text)

	newStatus := PalindromeWorkStatus{
		isPalindrome: isPalindrome,
		done:         true,
	}

	// pretend this is really slow
	var delay float64 = 0
	delay_str := os.Getenv("S_DELAY")
	if v, err := strconv.Atoi(delay_str); err == nil && v > 0 {
		delay = float64(v)
	}

	quarter_delay := delay / 4
	for i := 0; i < 4; i++ {
		time.Sleep(time.Duration(quarter_delay) * time.Second)

		p.lock.Lock()
		work, ok := p.work[msg.hash]
		p.lock.Unlock()

		if !ok {
			return
		} else {
			select {
			case <-work.cancel:
				return
			default:
			}
		}
	}

	// update work
	p.lock.Lock()
	if work, ok := p.work[msg.hash]; ok {
		work.status = newStatus
		p.work[msg.hash] = work
		for _, listener := range work.messages {
			select {
			case listener <- newStatus:
			default:
			}
		}
	}
	p.lock.Unlock()
}
