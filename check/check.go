package check

import (
	"strings"
)

func ContainsCorrectErrHandler(line string) bool {
	words := strings.Split(line, " ")
	returnHandler := -1
	ifHandler := -1
	hasErrBefore := false
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "return" {
			returnHandler = 1
		}
		if w == "if" {
			ifHandler = 3
		}
		if w == "err" {
			returnHandler--
			if !hasErrBefore {
				ifHandler--
			}
			hasErrBefore = true
		}
		if w == "!=" {
			ifHandler--
		}
		if w == "nil" {
			ifHandler--
		}
	}
	if strings.Contains(line, "if err != nil && strings.Contains(") {
		// println("line: ", line, fmt.Sprintf(":> %+v", ifHandler))
	}

	// println()
	return returnHandler == 0 || ifHandler == 0
}
