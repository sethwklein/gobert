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

func readSmallInt(buf *bytes.Buffer) int	{ return read1(buf) }

func readInt(buf *bytes.Buffer) int	{ return read4(buf) }

func readAtom(buf *bytes.Buffer) string {
	// Go doesn't have atom so treat them like strings
	return readString(buf)
}

func readSmallTuple(buf *bytes.Buffer) Term {
	size := read1(buf);
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
	size := int64(read2(buf));
	str, _ := ioutil.ReadAll(io.LimitReader(buf, size));
	return string(str);
}

func readList(buf *bytes.Buffer) []Term {
	size := read4(buf);
	list := make([]Term, size);

	for i := 0; i < size; i++ {
		list[i] = readTag(buf)
	}

	buf.ReadByte();

	return list;
}

func readBin(buf *bytes.Buffer) string {
	size := int64(read4(buf));
	str, _ := ioutil.ReadAll(io.LimitReader(buf, size));
	return string(str);
}

func readComplex(buf *bytes.Buffer) Term {
	if bytes.HasPrefix(buf.Bytes(), ComplexNil) {
		ioutil.ReadAll(io.LimitReader(buf, int64(len(ComplexNil))));
		return nil;
	}

	if bytes.HasPrefix(buf.Bytes(), ComplexTrue) {
		ioutil.ReadAll(io.LimitReader(buf, int64(len(ComplexTrue))));
		return true;
	}

	if bytes.HasPrefix(buf.Bytes(), ComplexFalse) {
		ioutil.ReadAll(io.LimitReader(buf, int64(len(ComplexFalse))));
		return false;
	}

	return nil;
}

func readTag(buf *bytes.Buffer) Term {
	tag, _ := buf.ReadByte();
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
	buf := bytes.NewBuffer(data);

	version, _ := buf.ReadByte();

	// check protocol version
	if version != Version {
		// Bad magic
		return -1
	}

	return readTag(buf);
}
