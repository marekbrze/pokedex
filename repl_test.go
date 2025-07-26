package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "hello #$%#$^@#%@# world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "12312423513424245234",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Input slices are not the same length: %v (%v) vs %v (%v)", actual, len(actual), c.expected, len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expected := c.expected[i]
			if word != expected {
				t.Errorf("Words are not the same: %v vs %v", word, expected)
			}
		}

	}
}
