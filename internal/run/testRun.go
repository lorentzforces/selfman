package run

import "testing"

func BailIfFailed(t *testing.T) {
	if t.Failed() {
		t.FailNow()
	}
}
