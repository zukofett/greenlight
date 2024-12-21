package assert

import (
	"slices"
	"strings"
	"testing"
)

func Equal[T comparable](t testing.TB, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t testing.TB, actual, expectedSubstring string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}

func NilError(t testing.TB, actual error) {
	t.Helper()

	if actual != nil {
		t.Errorf("got: %v; expected: nil", actual)
	}
}

func Contains[S ~[]T, T comparable](t testing.TB, actual S, expected T) {
	t.Helper()

	if !slices.Contains(actual, expected) {
		t.Errorf("got: %v; expected to have %v", actual, expected)
	}
}
