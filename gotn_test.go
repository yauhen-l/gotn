package main

import "testing"

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

func Test(t *testing.T) {
	t.Run("top level", func(t *testing.T) {
		expected := "Test"

		for _, pos := range []int{31, 57} {
			assertAtPos(t, expected, pos)
		}

		for _, pos := range []int{152, 153} {
			assertAtPos(t, expected, pos)
		}
	})

	t.Run("second level", func(t *testing.T) {
		expected := "Test/some_subset"

		for _, pos := range []int{58, 101} {
			assertAtPos(t, expected, pos)
		}

		for _, pos := range []int{148, 151} {
			assertAtPos(t, expected, pos)
		}
	})

	t.Run("third level", func(t *testing.T) {
		expected := "Test/some_subset/and_deeper"

		for _, pos := range []int{102, 147} {
			assertAtPos(t, expected, pos)
		}
	})

}

func assertAtPos(t *testing.T, expected string, pos int) {
	actual := getTestNameAtPos("test_data", []byte(testFileSrc), pos)
	if expected != actual {
		t.Errorf("pos %d - expected test name: %q, but got %q. Near: %q", pos, expected, actual, testFileSrc[pos:])
	}
}
