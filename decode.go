// http://bert-rpc.org/
// http://golang.org/

package bert

import (
	"bytes";
	"encoding/binary";
	"io";
	"io/ioutil";
	"os";
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

type Error struct {
	os.ErrorString;
}

var ErrBadMagic os.Error = &Error{"bad magic"}

type Term interface{}

func read1(buf *bytes.Buffer) (int, os.Error) {
	ui4, err := buf.ReadByte();
	if err != nil {
		return 0, err
	}

	return int(ui4), nil;
}

func read2(buf *bytes.Buffer) (int, os.Error) {
	bits, err := ioutil.ReadAll(io.LimitReader(buf, 2));
	if err != nil {
		return 0, err
	}

	ui16 := binary.BigEndian.Uint16(bits);
	return int(ui16), nil;
}

func read4(buf *bytes.Buffer) (int, os.Error) {
	bits, err := ioutil.ReadAll(io.LimitReader(buf, 4));
	if err != nil {
		return 0, err
	}

	ui32 := binary.BigEndian.Uint32(bits);
	return int(ui32), nil;
}

func readSmallInt(buf *bytes.Buffer) (int, os.Error) {
	return read1(buf)
}

func readInt(buf *bytes.Buffer) (int, os.Error) {
	return read4(buf)
}

func readAtom(buf *bytes.Buffer) (string, os.Error) {
	// Go doesn't have atom so treat them like strings
	return readString(buf)
}

func readSmallTuple(buf *bytes.Buffer) (Term, os.Error) {
	size, err := read1(buf);
	if err != nil {
		return nil, err
	}

	tuple := make([]Term, size);

	if bytes.HasPrefix(buf.Bytes(), ComplexBert) {
		return readComplex(buf), nil
	}

	for i := 0; i < size; i++ {
		term, err := readTag(buf);
		if err != nil {
			return nil, err
		}
		tuple[i] = term;
	}

	return tuple, nil;
}

func readNil(buf *bytes.Buffer) ([]Term, os.Error) {
	buf.ReadByte();
	list := make([]Term, 0);
	return list, nil;
}

func readString(buf *bytes.Buffer) (string, os.Error) {
	size, err := read2(buf);
	if err != nil {
		return "", err
	}

	str, err := ioutil.ReadAll(io.LimitReader(buf, int64(size)));
	if err != nil {
		return "", err
	}

	return string(str), nil;
}

func readList(buf *bytes.Buffer) ([]Term, os.Error) {
	size, err := read4(buf);
	if err != nil {
		return nil, err
	}

	list := make([]Term, size);

	for i := 0; i < size; i++ {
		term, err := readTag(buf);
		if err != nil {
			return nil, err
		}
		list[i] = term;
	}

	buf.ReadByte();

	return list, nil;
}

func readBin(buf *bytes.Buffer) (string, os.Error) {
	size, err := read4(buf);
	if err != nil {
		return "", err
	}

	str, err := ioutil.ReadAll(io.LimitReader(buf, int64(size)));
	if err != nil {
		return "", err
	}

	return string(str), nil;
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

func readTag(buf *bytes.Buffer) (Term, os.Error) {
	tag, err := buf.ReadByte();
	if err != nil {
		return nil, err
	}

	switch tag {
	case SmallInt:
		return readSmallInt(buf)
	case Int:
		return readInt(buf)
	case SmallBignum:
		return -1, nil
	case LargeBignum:
		return -1, nil
	case Float:
		return -1, nil
	case Atom:
		return readAtom(buf)
	case SmallTuple:
		return readSmallTuple(buf)
	case LargeTuple:
		return -1, nil
	case Nil:
		return readNil(buf)
	case String:
		return readString(buf)
	case List:
		return readList(buf)
	case Bin:
		return readBin(buf)
	}

	return nil, nil;
}

func Decode(data []byte) (Term, os.Error) {
	buf := bytes.NewBuffer(data);

	version, err := buf.ReadByte();

	if err != nil {
		return nil, err
	}

	// check protocol version
	if version != Version {
		return nil, ErrBadMagic
	}

	return readTag(buf);
}
