package bert

import (
	"testing";
	"reflect";
)

type Result struct {
	Name string;
}

func TestUnmarshal(t *testing.T) {
	var a struct {
		First string;
	}
	Unmarshal([]byte{131, 104, 1, 100, 0, 3, 102, 111, 111}, &a);
	assertEqual(t, "foo", a.First);

	var b struct {
		First int;
	}
	Unmarshal([]byte{131, 104, 1, 97, 42}, &b);
	assertEqual(t, 42, b.First);

	var c struct {
		First	string;
		Second	string;
	}
	Unmarshal([]byte{131, 104, 2, 100, 0, 3, 102, 111, 111, 100, 0, 3, 98, 97, 114}, &c);
	assertEqual(t, "foo", c.First);
	assertEqual(t, "bar", c.Second);

}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
