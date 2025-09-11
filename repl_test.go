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
			input:    "  hello   world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "THIS IS SPARTA",
			expected: []string{"this", "is", "sparta"},
		},
		{
			input:    "H e l l o M o t o ! ",
			expected: []string{"h", "e", "l", "l", "o", "m", "o", "t", "o", "!"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("%d != %d. Length of result does not match length of expected.", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%s != %s. Word does not match expected word", word, expectedWord)
			}
		}
	}
}
