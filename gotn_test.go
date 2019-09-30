package main

import (
	ts "testing"
)

const testFileSrc = `package main
import "testing"
func Test(t *testing.T) {
	t.Run("some subset", func(t *testing.T) {
		t.Run("and deeper", func(t *testing.T) {
		})
	})
}
func x(a int) {}
`

func TestGetTestNameAtPos(t *ts.T) {
	cases := []struct {
		expected string
		pos      int
	}{
		{"Test", 31},
		{"Test", 57},
		{"Test", 152},
		{"Test", 153},
		{"Test/some_subset", 58},
		{"Test/some_subset", 101},
		{"Test/some_subset", 148},
		{"Test/some_subset", 151},
		{"Test/some_subset/and_deeper", 102},
		{"Test/some_subset/and_deeper", 147},
	}
	for _, c := range cases {
		assertAtPos(t, c.expected, c.pos)
	}
}

func assertAtPos(t *ts.T, expected string, pos int) {
	actual := getTestNameAtPos("test_data", []byte(testFileSrc), pos)
	if expected != actual {
		t.Errorf("pos %d - expected test name: %q, but got %q. Near: %q", pos, expected, actual, testFileSrc[pos:])
	}
}
