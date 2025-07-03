package ffprocess

import (
	"testing"
)

func TestGetFloatString(t *testing.T) {
	val := 1.517
	want := "1.517"
	is := GetFloatString(val)
	if is != want {
		t.Errorf("GetFloatString:: Equal: is != want. %s != %s", is, want)
	}

	// longer prec
	val = 1.517222
	is = GetFloatString(val)
	if is != want {
		t.Errorf("GetFloatString:: Long prec: is != want. %s != %s", is, want)
	}
}

func TestGetMilliSeconds(t *testing.T) {
	val := 1.517
	want := 1517
	is := GetMilliSeconds(val)
	if is != want {
		t.Errorf("GetMilliSeconds:: Equal: is != want. %d != %d", is, want)
	}

	// longer prec
	val = 1.517222
	is = GetMilliSeconds(val)
	if is != want {
		t.Errorf("GetMilliSeconds:: Long prec: is != want. %d != %d", is, want)
	}
}
