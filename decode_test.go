package bert

import (
	"testing";
	"reflect";
)

func TestDecode(t *testing.T) {
	// Small Integer
	assertDecode(t, []byte{131, 97, 1}, 1);
	assertDecode(t, []byte{131, 97, 2}, 2);
	assertDecode(t, []byte{131, 97, 3}, 3);
	assertDecode(t, []byte{131, 97, 4}, 4);

	// Integer
	assertDecode(t, []byte{131, 98, 0, 0, 1, 1}, 257);
	assertDecode(t, []byte{131, 98, 0, 0, 4, 1}, 1025);
	assertDecode(t, []byte{131, 98, 255, 255, 255, 255}, -1);
	assertDecode(t, []byte{131, 98, 255, 255, 255, 248}, -8);

	// Small Bignum
	// assertDecode(t, []byte{131, 110, 8, 0, 0, 0, 232, 137, 4, 35, 199, 138}, 1);

	// Large Bignum

	// Float

	// Atom
	assertDecode(t, []byte{131, 100, 0, 3, 102, 111, 111}, "foo");

	// Small Tuple
	assertDecode(t, []byte{131, 104, 0}, []Term{});
	assertDecode(t, []byte{131, 104, 1, 100, 0, 3, 102, 111, 111}, []Term{"foo"});
	assertDecode(t, []byte{131, 104, 2, 100, 0, 3, 102, 111, 111, 100, 0, 3, 98, 97, 114}, []Term{"foo", "bar"});
	assertDecode(t, []byte{131, 104, 3, 100, 0, 5, 99, 111, 111, 114, 100, 97, 23, 97, 42}, []Term{"coord", 23, 42});

	// Large Tuple

	// String
	assertDecode(t, []byte{131, 107, 0, 3, 102, 111, 111}, "foo");

	// List
	assertDecode(t, []byte{131, 106}, []Term{});
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 97, 1, 106}, []Term{1});
	assertDecode(t, []byte{131, 108, 0, 0, 0, 3, 97, 1, 97, 2, 97, 3, 106}, []Term{1, 2, 3});
	//assertDecode(t, []byte{131, 108, 0, 0, 0, 2, 100, 0, 1, 97, 107, 0, 2, 1, 2, 106}, []Term{"a", []Term{1, 2}});

	// Binary
	assertDecode(t, []byte{131, 109, 0, 0, 0, 3, 102, 111, 111}, "foo");

	// Complex
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 3, 110, 105, 108}, nil);
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 4, 116, 114, 117, 101}, true);
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 5, 102, 97, 108, 115, 101}, false);
}

func assertDecode(t *testing.T, data []byte, expected interface{}) {
	val := Decode(data);
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", data, val, expected)
	}
}
