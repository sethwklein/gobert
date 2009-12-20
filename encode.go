package bert

import (
	"bytes";
	"os";
	"reflect";
	"fmt";
)

func write1(buf *bytes.Buffer, ui4 uint8)	{ buf.WriteByte(ui4) }

func writeSmallInt(buf *bytes.Buffer, n int) {
	write1(buf, SmallIntTag);
	write1(buf, uint8(n));
}

func writeTag(buf *bytes.Buffer, val interface{}) (err os.Error) {
	switch v := reflect.NewValue(val).(type) {
	case *reflect.IntValue:
		writeSmallInt(buf, v.Get())
	default:
		// TODO: Remove debug line
		fmt.Printf("Couldn't encode: %#v\n", v);
		err = ErrUnknownType;
	}

	return;
}

func Encode(val interface{}) ([]byte, os.Error) {
	buf := bytes.NewBuffer([]byte{});
	write1(buf, VersionTag);
	err := writeTag(buf, val);
	return buf.Bytes(), err;
}
