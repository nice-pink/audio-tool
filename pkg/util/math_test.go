package util

import (
	"testing"
)

func TestUnsynchSafe(t *testing.T) {
	var val uint32 = 1594
	var expect uint32 = 826
	res := Unsynchsafe(val)
	if res != expect {
		t.Error("res:", res, "expect:", expect)
	}
}

func TestUnsynchSafeSynchSafe(t *testing.T) {
	var val uint32 = 1023
	synch := Synchsafe(val)
	unsynch := Unsynchsafe(synch)
	if val != unsynch {
		t.Error("val:", val, "unsynch:", unsynch, "synch", synch)
	}
}
