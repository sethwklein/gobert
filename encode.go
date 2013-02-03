package bert

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func write1(w io.Writer, ui8 uint8) { w.Write([]byte{ui8}) }

func write2(w io.Writer, ui16 uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, ui16)
	w.Write(b)
}

func write4(w io.Writer, ui32 uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ui32)
	w.Write(b)
}

func writeSmallInt(w io.Writer, n uint8) {
	write1(w, SmallIntTag)
	write1(w, n)
}

func writeInt(w io.Writer, n uint32) {
	write1(w, IntTag)
	write4(w, n)
}

func writeFloat(w io.Writer, f float32) {
	write1(w, FloatTag)

	s := fmt.Sprintf("%.20e", float32(f))
	w.Write([]byte(s))

	pad := make([]byte, 31-len(s))
	w.Write(pad)
}

func writeAtom(w io.Writer, a string) {
	write1(w, AtomTag)
	write2(w, uint16(len(a)))
	w.Write([]byte(a))
}

func writeSmallTuple(w io.Writer, t reflect.Value) {
	write1(w, SmallTupleTag)
	size := t.Len()
	write1(w, uint8(size))

	for i := 0; i < size; i++ {
		writeTag(w, t.Index(i))
	}
}

func writeNil(w io.Writer) { write1(w, NilTag) }

func writeString(w io.Writer, s string) {
	write1(w, StringTag)
	write2(w, uint16(len(s)))
	w.Write([]byte(s))
}

func writeList(w io.Writer, l reflect.Value) {
	write1(w, ListTag)
	size := l.Len()
	write4(w, uint32(size))

	for i := 0; i < size; i++ {
		writeTag(w, l.Index(i))
	}

	writeNil(w)
}

func writeTag(w io.Writer, val reflect.Value) (err error) {
	switch v := val; v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := v.Int()
		if n >= 0 && n < 256 {
			writeSmallInt(w, uint8(n))
		} else {
			writeInt(w, uint32(n))
		}
	case reflect.Float32, reflect.Float64:
		writeFloat(w, float32(v.Float()))
	case reflect.String:
		if v.Type().Name() == "Atom" {
			writeAtom(w, v.String())
		} else {
			writeString(w, v.String())
		}
	case reflect.Slice:
		writeSmallTuple(w, v)
	case reflect.Array:
		writeList(w, v)
	case reflect.Interface:
		writeTag(w, v.Elem())
	default:
		if !reflect.Indirect(val).IsValid() {
			writeNil(w)
		} else {
			err = ErrUnknownType
		}
	}

	return
}

func EncodeTo(w io.Writer, val interface{}) (err error) {
	write1(w, VersionTag)
	err = writeTag(w, reflect.ValueOf(val))
	return
}

func Encode(val interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := EncodeTo(buf, val)
	return buf.Bytes(), err
}

func Marshal(w io.Writer, val interface{}) error {
	return EncodeTo(w, val)
}

func MarshalResponse(w io.Writer, val interface{}) (err error) {
	resp, err := Encode(val)

	write4(w, uint32(len(resp)))
	w.Write(resp)

	return
}
