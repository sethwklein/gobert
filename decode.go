// http://bert-rpc.org/
// http://golang.org/

package bert

import (
	"bytes";
	"encoding/binary";
	"io";
	"io/ioutil";
)

const (
	Version		= 131;
	SmallInt	= 97;
	Int		= 98;
	SmallBignum	= 110;
	LargeBignum	= 111;
	Float		= 99;
	Atom		= 100;
	SmallTuple	= 104;
	LargeTuple	= 105;
	Nil		= 106;
	String		= 107;
	List		= 108;
	Bin		= 109;
)

var (
	ComplexBert	= []uint8{100, 0, 4, 98, 101, 114, 116, 100};
	ComplexNil	= []uint8{100, 0, 4, 98, 101, 114, 116, 100, 0, 3, 110, 105, 108};
	ComplexTrue	= []uint8{100, 0, 4, 98, 101, 114, 116, 100, 0, 4, 116, 114, 117, 101};
	ComplexFalse	= []uint8{100, 0, 4, 98, 101, 114, 116, 100, 0, 5, 102, 97, 108, 115, 101};
)

type Term interface{}

func read1(buf *bytes.Buffer) int {
	ui4, _ := buf.ReadByte();
	return int(ui4);
}

func read2(buf *bytes.Buffer) int {
	bits, _ := ioutil.ReadAll(io.LimitReader(buf, 2));
        ui16 := binary.BigEndian.Uint16(bits);
        return int(ui16);
}

func read4(buf *bytes.Buffer) int {
	bits, _ := ioutil.ReadAll(io.LimitReader(buf, 4));
        ui32 := binary.BigEndian.Uint32(bits);
        return int(ui32);
}

func readSmallInt(buf *bytes.Buffer) int {
	return read1(buf);
}

func readInt(buf *bytes.Buffer) int {
	return read4(buf)
}

func readAtom(buf *bytes.Buffer) string {
	var size = read2(buf);
	var str = buf.Bytes()[0:size];
	for i := 0; i < size; i++ {
		buf.ReadByte()
	}
	return string(str);
}

func readSmallTuple(buf *bytes.Buffer) Term {
	var size = read1(buf);
	tuple := make([]Term, size);

	if bytes.HasPrefix(buf.Bytes(), ComplexBert) {
		return readComplex(buf)
	}

	for i := 0; i < size; i++ {
		tuple[i] = readTag(buf)
	}

	return tuple;
}

func readNil(buf *bytes.Buffer) []Term {
	buf.ReadByte();
	list := make([]Term, 0);
	return list;
}

func readString(buf *bytes.Buffer) string {
	var size = read2(buf);
	var str = buf.Bytes()[0:size];
	for i := 0; i < size; i++ {
		buf.ReadByte()
	}
	return string(str);
}

func readList(buf *bytes.Buffer) []Term {
	var size = read4(buf);
	list := make([]Term, size);

	for i := 0; i < size; i++ {
		list[i] = readTag(buf)
	}

	buf.ReadByte();

	return list;
}

func readBin(buf *bytes.Buffer) string {
	var size = read4(buf);
	var str = buf.Bytes()[0:size];
	return string(str);
}

func readComplex(buf *bytes.Buffer) Term {
	if bytes.HasPrefix(buf.Bytes(), ComplexNil) {
		for i := 0; i < len(ComplexNil); i++ {
			buf.ReadByte()
		}
		return nil;
	}

	if bytes.HasPrefix(buf.Bytes(), ComplexTrue) {
		for i := 0; i < len(ComplexTrue); i++ {
			buf.ReadByte()
		}
		return true;
	}

	if bytes.HasPrefix(buf.Bytes(), ComplexFalse) {
		for i := 0; i < len(ComplexFalse); i++ {
			buf.ReadByte()
		}
		return false;
	}

	return nil;
}

func readTag(buf *bytes.Buffer) Term {
	var tag, _ = buf.ReadByte();
	switch tag {
	case SmallInt:
		return readSmallInt(buf)
	case Int:
		return readInt(buf)
	case SmallBignum:
		return -1
	case LargeBignum:
		return -1
	case Float:
		return -1
	case Atom:
		return readAtom(buf)
	case SmallTuple:
		return readSmallTuple(buf)
	case LargeTuple:
		return -1
	case Nil:
		return readNil(buf)
	case String:
		return readString(buf)
	case List:
		return readList(buf)
	case Bin:
		return readBin(buf)
	}

	return -1;
}

func Decode(data []byte) Term {
	var buf = bytes.NewBuffer(data);

	var version, _ = buf.ReadByte();

	// check protocol version
	if version != Version {
		// Bad magic
		return -1
	}

	return readTag(buf);
}
