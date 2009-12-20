package bert

import (
	"testing";
	"reflect";
)

type Request struct {
	Kind		Atom;
	Module		Atom;
	Function	Atom;
	Arguments	[]Term;
}

type Response struct {
	Kind	string;
	Result	string;
}

func TestUnmarshal(t *testing.T) {
	var a struct {
		First Atom;
	}
	Unmarshal([]byte{131, 104, 1, 100, 0, 3, 102, 111, 111}, &a);
	assertEqual(t, Atom("foo"), a.First);

	var b struct {
		First int;
	}
	Unmarshal([]byte{131, 104, 1, 97, 42}, &b);
	assertEqual(t, 42, b.First);

	var c struct {
		First	Atom;
		Second	Atom;
	}
	Unmarshal([]byte{131, 104, 2, 100, 0, 3, 102, 111, 111, 100, 0, 3, 98, 97, 114}, &c);
	assertEqual(t, Atom("foo"), c.First);
	assertEqual(t, Atom("bar"), c.Second);

	var req Request;
	Unmarshal([]byte{131, 104, 4,
		100, 0, 4, 99, 97, 108, 108,
		100, 0, 6, 112, 104, 111, 116, 111, 120,
		100, 0, 8, 105, 109, 103, 95, 115, 105, 122, 101,
		108, 0, 0, 0, 1, 97, 99,
		106,
	},
		&req);
	assertEqual(t, Atom("call"), req.Kind);
	assertEqual(t, Atom("photox"), req.Module);
	assertEqual(t, Atom("img_size"), req.Function);
	assertEqual(t, []Term{99}, req.Arguments);
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
