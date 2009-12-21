package bert

import (
	"bytes";
	"encoding/binary";
	"fmt";
	"io";
	"os";
	"reflect";
	"strings";
)

func write1(w io.Writer, ui4 uint8)	{ w.Write([]byte{ui4}) }

func write2(w io.Writer, ui16 uint16) {
	b := make([]byte, 2);
	binary.BigEndian.PutUint16(b, ui16);
	w.Write(b);
}

func write4(w io.Writer, ui32 uint32) {
	b := make([]byte, 4);
	binary.BigEndian.PutUint32(b, ui32);
	w.Write(b);
}

func writeSmallInt(w io.Writer, n int) {
	write1(w, SmallIntTag);
	write1(w, uint8(n));
}

func writeInt(w io.Writer, n int) {
	write1(w, IntTag);
	write4(w, uint32(n));
}

func writeAtom(w io.Writer, a string) {
	write1(w, AtomTag);
	write2(w, uint16(len(a)));
	w.Write(strings.Bytes(a));
}

func writeSmallTuple(w io.Writer, t *reflect.SliceValue) {
	write1(w, SmallTupleTag);
	size := t.Len();
	write1(w, uint8(size));

	for i := 0; i < size; i++ {
		writeTag(w, t.Elem(i))
	}
}

func writeNil(w io.Writer)	{ write1(w, NilTag) }

func writeString(w io.Writer, s string) {
	write1(w, StringTag);
	write2(w, uint16(len(s)));
	w.Write(strings.Bytes(s));
}

func writeList(w io.Writer, l *reflect.ArrayValue) {
	write1(w, ListTag);
	size := l.Len();
	write4(w, uint32(size));

	for i := 0; i < size; i++ {
		writeTag(w, l.Elem(i))
	}

	writeNil(w);
}

func writeTag(w io.Writer, val reflect.Value) (err os.Error) {
	switch v := val.(type) {
	case *reflect.IntValue:
		n := v.Get();
		if n >= 0 && n < 256 {
			writeSmallInt(w, n)
		} else {
			writeInt(w, n)
		}
	case *reflect.StringValue:
		if v.Type().Name() == "Atom" {
			writeAtom(w, v.Get())
		} else {
			writeString(w, v.Get())
		}
	case *reflect.SliceValue:
		writeSmallTuple(w, v)
	case *reflect.ArrayValue:
		writeList(w, v)
	case *reflect.InterfaceValue:
		writeTag(w, v.Elem())
	default:
		if reflect.Indirect(val) == nil {
			writeNil(w)
		} else {
			// TODO: Remove debug line
			fmt.Printf("Couldn't encode: %#v\n", v);
			err = ErrUnknownType;
		}
	}

	return;
}

func Encode(val interface{}) ([]byte, os.Error) {
	buf := bytes.NewBuffer([]byte{});
	write1(buf, VersionTag);
	err := writeTag(buf, reflect.NewValue(val));
	return buf.Bytes(), err;
}
