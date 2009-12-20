package bert

import (
	"bytes";
	"encoding/binary";
	"os";
	"reflect";
	"fmt";
)

func write1(buf *bytes.Buffer, ui4 uint8)	{ buf.WriteByte(ui4) }

func write2(buf *bytes.Buffer, ui16 uint16) {
	b := make([]byte, 2);
	binary.BigEndian.PutUint16(b, ui16);
	buf.Write(b);
}

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

func writeAtom(buf *bytes.Buffer, a string) {
	write1(buf, AtomTag);
	write2(buf, uint16(len(a)));
	buf.WriteString(a);
}

func writeSmallTuple(buf *bytes.Buffer, t *reflect.SliceValue) {
	write1(buf, SmallTupleTag);
	size := t.Len();
	write1(buf, uint8(size));

	for i := 0; i < size; i++ {
		writeTag(buf, t.Elem(i))
	}
}

func writeTag(buf *bytes.Buffer, val reflect.Value) (err os.Error) {
	switch v := val.(type) {
	case *reflect.IntValue:
		n := v.Get();
		if n >= 0 && n < 256 {
			writeSmallInt(buf, n)
		} else {
			writeInt(buf, n)
		}
	case *reflect.StringValue:
		if v.Type().Name() == "Atom" {
			writeAtom(buf, v.Get())
		} else {
			err = ErrUnknownType
		}
	case *reflect.SliceValue:
		writeSmallTuple(buf, v)
	case *reflect.InterfaceValue:
		writeTag(buf, v.Elem())
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
	err := writeTag(buf, reflect.NewValue(val));
	return buf.Bytes(), err;
}
