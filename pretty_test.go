package pretty

import (
	"fmt"
	"testing"
)

func doesntPanic(t *testing.T) {
	if err := recover(); err != nil {
		t.Fatal(fmt.Sprint(err))
	}
}

func TestString(t *testing.T) {
	defer doesntPanic(t)

	tests := []string{
		"foo bar }",
		"foo bar {",
		"foo bar ,",
		"foo bar , baz",
		"foo bar } baz",
		"foo bar { baz",
		"foo bar {{}baz",
	}
	for _, test := range tests {
		if got := Print(test); got != fmt.Sprintf("%q", test) {
			t.Error(fmt.Sprintf("%q != %q", fmt.Sprintf("%q", test), got))
		}
	}
}
