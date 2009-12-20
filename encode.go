package bert

import (
	"bytes";
	"encoding/binary";
	"os";
	"reflect";
	"fmt";
)

func write1(buf *bytes.Buffer, ui4 uint8)	{ buf.WriteByte(ui4) }

func write4(buf *bytes.Buffer, ui32 uint32) {
	b := make([]byte, 4);
	binary.BigEndian.PutUint32(b, ui32);
	buf.Write(b);
}

func writeSmallInt(buf *bytes.Buffer, n int) {
	write1(buf, SmallIntTag);
	write1(buf, uint8(n));
}

func writeInt(buf *bytes.Buffer, n int) {
	write1(buf, IntTag);
	write4(buf, uint32(n));
}

func writeTag(buf *bytes.Buffer, val interface{}) (err os.Error) {
	switch v := reflect.NewValue(val).(type) {
	case *reflect.IntValue:
		n := v.Get();
		if n >= 0 && n < 256 {
			writeSmallInt(buf, n)
		} else {
			writeInt(buf, n)
		}
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
