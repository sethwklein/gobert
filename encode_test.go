package bert

import (
	"testing";
	"reflect";
)

func TestEncode(t *testing.T) {
	assertEncode(t, 1, []byte{131, 97, 1});
	assertEncode(t, 42, []byte{131, 97, 42});
}

func assertEncode(t *testing.T, actual interface{}, expected []byte) {
	val, err := Encode(actual);
	if err != nil {
		t.Errorf("Encode(%v) returned error '%v'", actual, err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("Decode(%v) = %v, expected %v", actual, val, expected)
	}
}
