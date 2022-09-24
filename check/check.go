package check

import (
	"strings"
)

func ContainsCorrectErrHandler(line string) bool {
	words := strings.Split(line, " ")

	errVar := "err"

	returnHandler := -1
	ifHandler := -1
	isNotNil := -1

	hasIfHandler := false
	hasErrBefore := false
	hasNotEqBefore := false

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
			ifHandler--
			isNotNil--
		}
	}

	if hasIfHandler && ifHandler != 0 {
		return strings.Contains(strings.ToLower(line), "is(err)")
	}

	return returnHandler == 0 || ifHandler == 0 || isNotNil == 0
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
