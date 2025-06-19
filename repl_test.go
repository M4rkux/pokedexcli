package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   multiple    spaces   here   ",
			expected: []string{"multiple", "spaces", "here"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i, expectedWord := range c.expected {
			word := actual[i]

			if word != expectedWord {
				t.Errorf("Word: %s doesn't match the expected word: %s", word, expectedWord)
			}
		}
	}
}
