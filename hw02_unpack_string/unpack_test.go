package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "aaa0b", expected: "aab"},
		{input: "", expected: ""},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: `\65`, expected: `66666`},
		{input: `\6a`, expected: `6a`},
		{input: `\6ab3c`, expected: `6abbbc`},
		{input: "\n", expected: "\n"},
		{input: "\r\n", expected: "\r\n"},
		{input: "\n3", expected: "\n\n\n"},
		{input: `\\5`, expected: `\\\\\`},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: "\"3", expected: `"""`},
		{input: "\a\b\f\n\r\t\v\\\\'\"", expected: "\a\b\f\n\r\t\v\\'\""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackDifferentLanguagesString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "Отус!8\\52\n", expected: "Отус!!!!!!!!55\n"},
		{input: "奧3圖2斯奧0\\52\r\n", expected: "奧奧奧圖圖斯55\r\n"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `abc\`, `qw\ne`, `123\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
