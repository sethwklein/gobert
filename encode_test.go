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

	// Float
	assertEncode(t, 0.5, []byte{131, 99, 53, 46, 48, 48, 48, 48, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 101,
		45, 48, 49, 0, 0, 0, 0, 0,
	});
	assertEncode(t, 3.14159, []byte{131, 99, 51, 46, 49, 52, 49, 53, 57,
		48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48, 48,
		101, 43, 48, 48, 0, 0, 0, 0, 0,
	});
	assertEncode(t, -3.14159, []byte{131, 99, 45, 51, 46, 49, 52, 49, 53,
		57, 48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48,
		48, 101, 43, 48, 48, 0, 0, 0, 0,
	});

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

	// Nil
	assertEncode(t, nil, []byte{131, 106});

	// String
	assertEncode(t, "foo", []byte{131, 107, 0, 3, 102, 111, 111});

	// List
	assertEncode(t, [1]Term{1},
		[]byte{131, 108, 0, 0, 0, 1, 97, 1, 106});
	assertEncode(t, [3]Term{1, 2, 3},
		[]byte{131, 108, 0, 0, 0, 3,
			97, 1, 97, 2, 97, 3,
			106,
		});
	assertEncode(t, [2]Term{Atom("a"), Atom("b")},
		[]byte{131, 108, 0, 0, 0, 2,
			100, 0, 1, 97, 100, 0, 1, 98,
			106,
		});
}

func assertEncode(t *testing.T, actual interface{}, expected []byte) {
	val, err := Encode(actual);
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
