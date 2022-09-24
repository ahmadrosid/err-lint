package check_test

import (
	"err-lint/check"
	"testing"
)

func TestValidateContains(t *testing.T) {
	scenarios := []struct {
		input    string
		expected bool
	}{
		{
			input:    "if err != nil {",
			expected: true,
		},
		{
			input:    "return err",
			expected: true,
		},
		{
			input:    "}); err != nil {",
			expected: true,
		},
		{
			input:    "\tif err != nil && err != redis.ErrNil {",
			expected: true,
		},
		{
			input:    "\tif err != nil && strings.Contains(",
			expected: true,
		},
		{
			input:    "return fmt, err",
			expected: true,
		},
		{
			input:    "return (SomeStruct)(*detail).ToEntity(), nil",
			expected: false,
		},
		{
			input:    "CheckErr(err)",
			expected: true,
		},
		{
			input:    "if err != nil || gitTree == nil {",
			expected: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.input, func(t *testing.T) {
			actual := check.ContainsCorrectErrHandler(s.input)
			if actual != s.expected {
				t.Errorf("\ninput: '%s'\nexpected: %+v got %+v", s.input, s.expected, actual)
			}
		})
	}
}
