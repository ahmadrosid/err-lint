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

	returnHandler := -1
	ifHandler := -1
	errHandler := -1

	hasIfHandler := false
	hasErrBefore := false
	hasNotEqBefore := false
	hasEqBefore := false

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "return" {
			returnHandler = 1
		}
		if w == "if" {
			ifHandler = 2
			hasIfHandler = true
		}
		if w == "err" {
			returnHandler--
			if !hasErrBefore {
				ifHandler--
				errHandler = 1
			}
			hasErrBefore = true
		}
		if w == "!=" {
			if !hasNotEqBefore {
				ifHandler--
				errHandler--
			}
			hasNotEqBefore = true
		}
		if w == "==" {
			if !hasEqBefore && !hasNotEqBefore && !hasErrBefore {
				ifHandler--
				errHandler--
			}
			if hasErrBefore && !hasEqBefore {
				errHandler--
			}

			hasEqBefore = true
		}
	}

	status := returnHandler == 0 || ifHandler == 0 || errHandler == 0
	if hasIfHandler && ifHandler == 0 && !status {
		return isOkErrorHandler(line)
	}

	if strings.Contains(line, "err") && !status {
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
