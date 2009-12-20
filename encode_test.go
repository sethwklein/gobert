package bert

import (
	"testing";
	"reflect";
)

func TestEncode(t *testing.T) {
	// Small Integer
	assertEncode(t, 1, []byte{131, 97, 1});
	assertEncode(t, 42, []byte{131, 97, 42});

	// Integer
	assertEncode(t, 257, []byte{131, 98, 0, 0, 1, 1});
	assertEncode(t, 1025, []byte{131, 98, 0, 0, 4, 1});
	assertEncode(t, -1, []byte{131, 98, 255, 255, 255, 255});
	assertEncode(t, -8, []byte{131, 98, 255, 255, 255, 248});
	assertEncode(t, 5000, []byte{131, 98, 0, 0, 19, 136});
	assertEncode(t, -5000, []byte{131, 98, 255, 255, 236, 120});

	// Atom
	assertEncode(t, Atom("foo"),
		[]byte{131, 100, 0, 3, 102, 111, 111});
}

func assertEncode(t *testing.T, actual interface{}, expected []byte) {
	val, err := Encode(actual);
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
