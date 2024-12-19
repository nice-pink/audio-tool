package riff

import "testing"

func TestBuild(t *testing.T) {
	var size uint32 = 1023
	data := Build(size, 0)
	if len(data) == 0 {
		t.Error("no size")
	}

	if !IsValid(data) {
		t.Error("is not valid")
	}

	tagSize := GetTagSize(data)
	if tagSize != size {
		t.Error("tagSize != expectedSize", tagSize, size)
	}
}
