// http://bert-rpc.org/
// http://golang.org/

package bert

import (
	"bytes";
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

func readBytes(buf *bytes.Buffer, bytes int) int {
	var n = 0;
	var c byte;

	for b := uint8(bytes * 8); b > 0; {
		b -= 8;
		c, _ = buf.ReadByte();
		n += int(c) << b;
	}

	return n;
}


func readList(buf *bytes.Buffer) []Term {
	var size = readBytes(buf, 4);
	list := make([]Term, size);

	for i := 0; i < size; i++ {
		list[i] = readTag(buf)
	}

	buf.ReadByte();

	return list;
}

func readSmallInt(buf *bytes.Buffer) int	{ return readBytes(buf, 1) }

func readInt(buf *bytes.Buffer) int {
	var value = readBytes(buf, 4);
	return value;
}

func readAtom(buf *bytes.Buffer) string {
	var size = readBytes(buf, 2);
	var str = buf.Bytes()[0:size];
	for i := 0; i < size; i++ {
		buf.ReadByte()
	}
	return string(str);
}

func readSmallTuple(buf *bytes.Buffer) Term {
	var size = readBytes(buf, 1);
	tuple := make([]Term, size);

	if bytes.HasPrefix(buf.Bytes(), ComplexBert) {
		return readComplex(buf)
	}

	for i := 0; i < size; i++ {
		tuple[i] = readTag(buf)
	}

	return tuple;
}

func readString(buf *bytes.Buffer) string {
	var size = readBytes(buf, 2);
	var str = buf.Bytes()[0:size];
	for i := 0; i < size; i++ {
		buf.ReadByte()
	}
	return string(str);
}

func readNil(buf *bytes.Buffer) []Term {
	buf.ReadByte();
	list := make([]Term, 0);
	return list;
}

func readBin(buf *bytes.Buffer) string {
	var size = readBytes(buf, 4);
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
