package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"err-lint/check"
	"err-lint/stack"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

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

func max(cur int, limit int) int {
	if cur >= limit {
		return limit
	}
	return cur
}

func filterScope(i int, lines []string) (next int) {
	line := lines[i]
	length := len(lines)

	countBracket := stack.NewStack()

	if strings.Contains(line, "{") || strings.Contains(line, "(") {
		countBracket = check.Bracket(countBracket, line)
		for {
			i++
			if i == length-1 {
				break
			}
			next = i
			nextLine := lines[next]
			countBracket = check.Bracket(countBracket, nextLine)
			hasDotSuffix := false
			if strings.HasSuffix(strings.TrimSpace(nextLine), ".") {
				hasDotSuffix = true
				continue
			}

			if countBracket.Len() == 0 {
				lastLine := strings.TrimSpace(lines[max(next+1, length-1)])
				if hasDotSuffix {
					hasDotSuffix = false
					continue
				}

				if strings.HasPrefix(lastLine, "//") {
					continue
				}

				if check.IsOnlyWhiteSpace(lastLine) {
					continue
				}
				break
			}
		}
	}
	return next
}

func Detect(filename string) {
	if !strings.HasSuffix(filename, ".go") {
		return
	}

	if strings.HasSuffix(filename, "_test.go") {
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
		if check.ContainsCorrectErrHandler(line) {
			continue
		}

		if strings.Contains(line, "err :=") {
			next := i + 1
			if next < length {
				if check.ContainsCorrectErrHandler(lines[next]) {
					continue
				}

				if strings.Contains(lines[next], "return") && strings.Contains(lines[next], "err") {
					continue
				}

				i := filterScope(i, lines)
				next = i + 1
				nextLine := lines[next]

				if check.ContainsCorrectErrHandler(nextLine) {
					continue
				}

				printedLines := []string{
					yellow(fmt.Sprintf("%s:%d", filename, next+1)),
					green(fmt.Sprintf("%d %s", curLineIdx+1, lines[curLineIdx])),
				}

				rangeIdx := next - curLineIdx
				dontPrint := false
				for {
					curLineIdx++
					rangeIdx--
					if curLineIdx == length {
						break
					}
					rangeLine := fmt.Sprintf("%d %s", curLineIdx+1, lines[curLineIdx])
					if check.ContainsCorrectErrHandler(rangeLine) {
						dontPrint = true
						break
					}
					if rangeIdx == 0 {
						printedLines = append(printedLines, red(rangeLine))
						break
					} else {
						printedLines = append(printedLines, green(rangeLine))
					}
				}

				if dontPrint {
					continue
				}
				println(strings.Join(printedLines, "\n"))
				println()
			}
		}
	}
}

func main() {
	source := "."
	if len(os.Args) > 1 {
		source = os.Args[1]
	}

	err := ReadDirectory(source, Detect)
	if err != nil {
		Detect(source)
	}
}
