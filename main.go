package main

import (
	"err-lint/stack"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReadDirectory(dir string, result func(filename string)) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				result(path)
			}
			return nil
		})
	if err != nil {
		return err
	}

	return nil
}

func isWhitespace(c rune) bool {
	switch c {
	case ' ', '\t', '\n', '\u000b', '\u000c', '\r':
		return true
	}
	return false
}

func isOnlyWhiteSpace(line string) bool {
	for _, c := range line {
		if !isWhitespace(c) {
			return false
		}
	}
	return true
}

func checkBracket(st *stack.Stack, line string) *stack.Stack {
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

func filterScope(i int, lines []string) (next int, skip bool) {
	next = i + 1
	line := lines[i]
	length := len(lines)

	countBracket := stack.NewStack()

	if strings.Contains(line, "{") || strings.Contains(line, "(") {
		countBracket = checkBracket(countBracket, line)
		for {
			i++
			if i >= length {
				break
			}
			next = i
			nextLine := lines[next]
			countBracket = checkBracket(countBracket, nextLine)

			if countBracket.Len() == 0 {
				break
			}
		}
	}

	return next, false
}

func Detect(filename string) {
	if !strings.HasSuffix(filename, ".go") {
		return
	}

	res, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	lines := strings.Split(string(res), "\n")
	length := len(lines)
	for i := 0; i < length; i++ {
		curLineIdx := i
		line := lines[curLineIdx]
		if strings.Contains(line, "err != nil") {
			continue
		}

		if strings.Contains(line, "err :=") {
			next := i + 1
			if next < length {
				i, skip := filterScope(i, lines)
				if skip {
					break
				}
				nextLine := lines[i+1]
				next = i + 1
				if strings.Contains(nextLine, "err != nil") {
					continue
				}
				if strings.Contains(nextLine, "return err") {
					continue
				}

				if isOnlyWhiteSpace(nextLine) {
					i++
					if i >= length {
						break
					}
					next = i
				}

				println(fmt.Sprintf("%s:%d", filename, next+1))
				println(fmt.Sprintf("%d %s", curLineIdx+1, lines[curLineIdx]))

				rangeIdx := next - curLineIdx - 1
				for {
					curLineIdx++
					println(fmt.Sprintf("%d %s", curLineIdx+1, lines[curLineIdx]))
					rangeIdx--
					if rangeIdx == 0 {
						break
					}
				}

				println()
			}
		}
	}
}

func main() {
	var filename = *flag.String("source", ".", "Source file or directory name.")
	flag.Parse()

	err := ReadDirectory(filename, Detect)
	if err != nil {
		Detect(filename)
	}
}
