package internal

import "testing"

func TestExtractHostname(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "user@hostname:22",
			expected: "hostname",
		},
		{
			input:    "user@hostname",
			expected: "hostname",
		},
		{
			input:    "hostname:22",
			expected: "hostname",
		},
		{
			input:    "hostname",
			expected: "hostname",
		},
		{
			input:    "ssh://user@hostname:22",
			expected: "hostname",
		},
		{
			input:    "ssh://user@hostname",
			expected: "hostname",
		},
		{
			input:    "ssh://hostname:22",
			expected: "hostname",
		},
		{
			input:    "ssh://hostname",
			expected: "hostname",
		},
	}

	for _, testCase := range testCases {
		actual := extractHostname(testCase.input)
		if actual != testCase.expected {
			t.Errorf("unexpected result. expected: %s, but got: %s", testCase.expected, actual)
		}
	}
}
