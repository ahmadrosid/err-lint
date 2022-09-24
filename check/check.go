package check

import (
	"err-lint/stack"
	"strings"
)

func Bracket(st *stack.Stack, line string) *stack.Stack {
	pair := map[rune]rune{
		'{': '}',
		'(': ')',
	}
	line = strings.TrimSpace(line)
	for _, c := range line {
		if _, exists := pair[c]; exists {
			st.Push(c)
		} else {
			newStack, ch := st.Pop()
			if ch == '0' {
				continue
			}
			if pair[ch] == c {
				st = newStack
			}
		}
	}
	return st
}

func isOkErrorHandler(line string) bool {
	line = strings.ToLower(line)
	return strings.Contains(line, "is(err") || strings.Contains(line, "(err)")
}

func ContainsCorrectErrHandler(line string) bool {
	words := strings.Split(line, " ")

	errVar := "err"

	returnHandler := -1
	ifHandler := -1
	isNotNil := -1

	hasIfHandler := false
	hasErrBefore := false
	hasNotEqBefore := false
	hasNilBefore := false

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "return" {
			returnHandler = 1
		}
		if w == "if" {
			ifHandler = 3
			hasIfHandler = true
		}
		if strings.Contains(w, errVar) {
			returnHandler--
			if !hasErrBefore {
				ifHandler--
				isNotNil = 2
			}
			hasErrBefore = true
		}
		if w == "!=" {
			if !hasNotEqBefore {
				ifHandler--
				isNotNil--
			}
			hasNotEqBefore = true
		}
		if w == "nil" {
			if !hasNilBefore {
				ifHandler--
				isNotNil--
				hasNilBefore = true
			}
		}
	}

	status := returnHandler == 0 || ifHandler == 0 || isNotNil == 0

	if hasIfHandler && ifHandler != 0 {
		return isOkErrorHandler(line)
	}

	if hasErrBefore && !status {
		return isOkErrorHandler(line)
	}

	return status
}

func isWhitespace(c rune) bool {
	switch c {
	case ' ', '\t', '\n', '\u000b', '\u000c', '\r':
		return true
	}
	return false
}

func IsOnlyWhiteSpace(line string) bool {
	for _, c := range line {
		if !isWhitespace(c) {
			return false
		}
	}
	return true
}
