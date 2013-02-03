package bert

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	// Small Integer
	assertDecode(t, []byte{131, 97, 1}, 1)
	assertDecode(t, []byte{131, 97, 2}, 2)
	assertDecode(t, []byte{131, 97, 3}, 3)
	assertDecode(t, []byte{131, 97, 4}, 4)
	assertDecode(t, []byte{131, 97, 42}, 42)

	// Integer
	assertDecode(t, []byte{131, 98, 0, 0, 1, 1}, 257)
	assertDecode(t, []byte{131, 98, 0, 0, 4, 1}, 1025)
	assertDecode(t, []byte{131, 98, 255, 255, 255, 255}, -1)
	assertDecode(t, []byte{131, 98, 255, 255, 255, 248}, -8)
	assertDecode(t, []byte{131, 98, 0, 0, 19, 136}, 5000)
	assertDecode(t, []byte{131, 98, 255, 255, 236, 120}, -5000)

	// Small Bignum
	// assertDecode(t, []byte{131, 110, 4, 0, 177, 104, 222, 58},
	// 	987654321);
	// assertDecode(t, []byte{131, 110, 4, 1, 177, 104, 222, 58},
	// 	-987654321);

	// Large Bignum

	// Float
	assertDecode(t, []byte{131, 99, 53, 46, 48, 48, 48, 48, 48, 48, 48,
		48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 101, 45,
		48, 49, 0, 0, 0, 0, 0,
	},
		float32(0.5))
	assertDecode(t, []byte{131, 99, 51, 46, 49, 52, 49, 53, 56,
		57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 56, 56, 50, 54, 50,
		101, 43, 48, 48, 0, 0, 0, 0, 0,
	},
		float32(3.14159))
	assertDecode(t, []byte{131, 99, 45, 51, 46, 49, 52, 49, 53,
		56, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 56, 56, 50, 54,
		50, 101, 43, 48, 48, 0, 0, 0, 0,
	},
		float32(-3.14159))
	assertDecode(t, []byte{131, 99, 51, 46, 49, 52, 49, 53, 57,
		48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48, 48,
		101, 43, 48, 48, 0, 0, 0, 0, 0,
	},
		float32(3.14159))
	assertDecode(t, []byte{131, 99, 45, 51, 46, 49, 52, 49, 53,
		57, 48, 49, 49, 56, 52, 48, 56, 50, 48, 51, 49, 50, 53, 48,
		48, 101, 43, 48, 48, 0, 0, 0, 0,
	},
		float32(-3.14159))

	// Atom
	assertDecode(t, []byte{131, 100, 0, 3, 102, 111, 111},
		Atom("foo"))
	assertDecode(t, []byte{131, 100, 0, 5, 104, 101, 108, 108, 111},
		Atom("hello"))

	// Small Tuple
	assertDecode(t, []byte{131, 104, 0}, []Term{})
	assertDecode(t, []byte{131, 104, 1,
		100, 0, 3, 102, 111, 111,
	},
		[]Term{Atom("foo")})
	assertDecode(t, []byte{131, 104, 2,
		100, 0, 3, 102, 111, 111,
		100, 0, 3, 98, 97, 114,
	},
		[]Term{Atom("foo"), Atom("bar")})
	assertDecode(t, []byte{131, 104, 3,
		100, 0, 5, 99, 111, 111, 114, 100,
		97, 23,
		97, 42,
	},
		[]Term{Atom("coord"), 23, 42})
	assertDecode(t, []byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		[]Term{Atom("call"), Atom("photox"), Atom("img_size"), []Term{99}})

	// Large Tuple

	// String
	assertDecode(t, []byte{131, 107, 0, 3, 102, 111, 111}, "foo")
	assertDecode(t, []byte{131, 107, 0, 1, 0}, "\000")
	assertDecode(t, []byte{131, 107, 0, 1, 1}, "\001")

	// List
	assertDecode(t, []byte{131, 106}, []Term{})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 97, 1, 106},
		[]Term{1})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 98, 0, 0, 1, 0, 106},
		[]Term{256})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 107, 0, 1, 97, 106},
		[]Term{"a"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 100, 0, 1, 97, 106},
		[]Term{Atom("a")})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 3,
		97, 1, 97, 2, 97, 3,
		106,
	},
		[]Term{1, 2, 3})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		107, 0, 1, 97, 107, 0, 1, 98, 106,
	},
		[]Term{"a", "b"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 1, 107, 0, 2, 97, 98, 106},
		[]Term{"ab"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		100, 0, 1, 97, 100, 0, 1, 98, 106,
	},
		[]Term{Atom("a"), Atom("b")})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2, 100, 0, 1, 97, 97, 1, 106},
		[]Term{Atom("a"), 1})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		107, 0, 1, 97, 107, 0, 2, 1, 2, 106,
	},
		[]Term{"a", "\001\002"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2,
		100, 0, 1, 97, 107, 0, 2, 1, 2, 106,
	},
		[]Term{Atom("a"), "\001\002"})
	assertDecode(t, []byte{131, 108, 0, 0, 0, 2, 100, 0, 1, 97, 108, 0, 0, 0, 1, 98, 0, 0, 1, 0, 106,
		106,
	},
		[]Term{Atom("a"), []Term{256}})

	// Binary
	assertDecode(t, []byte{131, 109, 0, 0, 0, 3, 102, 111, 111},
		[]uint8{102, 111, 111})
	assertDecode(t, []byte{131, 109, 0, 0, 0, 5, 104, 101, 108, 108, 111},
		[]uint8{104, 101, 108, 108, 111})

	// Complex
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 3, 110, 105, 108}, nil)
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 4, 116, 114, 117, 101}, true)
	assertDecode(t, []byte{131, 104, 2, 100, 0, 4, 98, 101, 114, 116, 100, 0, 5, 102, 97, 108, 115, 101}, false)

	assertDecode(t, []byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		[]Term{Atom("call"), Atom("photox"), Atom("img_size"), []Term{99}})
}

func assertDecode(t *testing.T, data []byte, expected interface{}) {
	val, err := Decode(data)
	if err != nil {
		t.Errorf("Decode(%v) returned error '%v'", data, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %#v, expected %#v", data, val, expected)
	}
}

func TestUnmarshal(t *testing.T) {
	var a struct {
		First Atom
	}
	Unmarshal([]byte{131, 104, 1, 100, 0, 3, 102, 111, 111}, &a)
	assertEqual(t, Atom("foo"), a.First)

	var b struct {
		First int
	}
	Unmarshal([]byte{131, 104, 1, 97, 42}, &b)
	assertEqual(t, 42, b.First)

	var c struct {
		First  Atom
		Second Atom
	}
	Unmarshal([]byte{131, 104, 2, 100, 0, 3, 102, 111, 111, 100, 0, 3, 98, 97, 114}, &c)
	assertEqual(t, Atom("foo"), c.First)
	assertEqual(t, Atom("bar"), c.Second)

	var req Request
	Unmarshal([]byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		&req)
	assertEqual(t, Atom("call"), req.Kind)
	assertEqual(t, Atom("photox"), req.Module)
	assertEqual(t, Atom("img_size"), req.Function)
	assertEqual(t, []Term{99}, req.Arguments)
}

func TestUnmarshalRequest(t *testing.T) {
	buf := bytes.NewBuffer([]byte{
		0, 0, 0, 38,
		131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	})

	req, _ := UnmarshalRequest(buf)
	assertEqual(t, Atom("call"), req.Kind)
	assertEqual(t, Atom("photox"), req.Module)
	assertEqual(t, Atom("img_size"), req.Function)
	assertEqual(t, []Term{99}, req.Arguments)
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
