package test

import (
	"runtime/debug"
	"testing"
)

func assertEq(t *testing.T, expected, actual uint64) {
	if expected != actual {
		t.Errorf("%s assertEq failed: expected %+v, got %+v\n%s", t.Name(), expected, actual, debug.Stack())
	}
}
