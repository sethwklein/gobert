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

	// Small Tuple
	assertEncode(t, []Term{Atom("foo")},
		[]byte{131, 104, 1, 100, 0, 3, 102, 111, 111});
	assertEncode(t, []Term{Atom("foo"), Atom("bar")},
		[]byte{131, 104, 2,
			100, 0, 3, 102, 111, 111,
			100, 0, 3, 98, 97, 114,
		});
	assertEncode(t, []Term{Atom("coord"), 23, 42},
		[]byte{131, 104, 3,
			100, 0, 5, 99, 111, 111, 114, 100,
			97, 23,
			97, 42,
		});

	// String
	assertEncode(t, "foo", []byte{131, 107, 0, 3, 102, 111, 111});
}

func assertEncode(t *testing.T, actual interface{}, expected []byte) {
	val, err := Encode(actual);
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
