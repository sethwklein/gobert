package bert

import (
	"bytes";
	"encoding/binary";
	"io";
	"io/ioutil";
	"os";
)

var ErrBadMagic os.Error = &Error{"bad magic"}
var ErrUnknownType os.Error = &Error{"unknown type"}

func read1(r io.Reader) (int, os.Error) {
	bits, err := ioutil.ReadAll(io.LimitReader(r, 1));
	if err != nil {
		return 0, err
	}

	ui8 := uint8(bits[0]);
	return int(ui8), nil;
}

func read2(r io.Reader) (int, os.Error) {
	bits, err := ioutil.ReadAll(io.LimitReader(r, 2));
	if err != nil {
		return 0, err
	}

	ui16 := binary.BigEndian.Uint16(bits);
	return int(ui16), nil;
}

func read4(r io.Reader) (int, os.Error) {
	bits, err := ioutil.ReadAll(io.LimitReader(r, 4));
	if err != nil {
		return 0, err
	}

	ui32 := binary.BigEndian.Uint32(bits);
	return int(ui32), nil;
}

func readSmallInt(r io.Reader) (int, os.Error) {
	return read1(r)
}

func readInt(r io.Reader) (int, os.Error)	{ return read4(r) }

func readAtom(r io.Reader) (Atom, os.Error) {
	str, err := readString(r);
	return Atom(str), err;
}

var (
	ComplexBert = []uint8{100, 0, 4, 98, 101, 114, 116, 100};
)

func readSmallTuple(r io.Reader) (Term, os.Error) {
	size, err := read1(r);
	if err != nil {
		return nil, err
	}

	tuple := make([]Term, size);

	for i := 0; i < size; i++ {
		term, err := readTag(r);
		if err != nil {
			return nil, err
		}
		switch a := term.(type) {
		case Atom:
			if a == BertAtom {
				return readComplex(r)
			}
		}
		tuple[i] = term;
	}

	return tuple, nil;
}

func readNil(r io.Reader) ([]Term, os.Error) {
	_, err := ioutil.ReadAll(io.LimitReader(r, 1));
	if err != nil {
		return nil, err
	}
	list := make([]Term, 0);
	return list, nil;
}

func readString(r io.Reader) (string, os.Error) {
	size, err := read2(r);
	if err != nil {
		return "", err
	}

	str, err := ioutil.ReadAll(io.LimitReader(r, int64(size)));
	if err != nil {
		return "", err
	}

	return string(str), nil;
}

func readList(r io.Reader) ([]Term, os.Error) {
	size, err := read4(r);
	if err != nil {
		return nil, err
	}

	list := make([]Term, size);

	for i := 0; i < size; i++ {
		term, err := readTag(r);
		if err != nil {
			return nil, err
		}
		list[i] = term;
	}

	read1(r);

	return list, nil;
}

func readBin(r io.Reader) ([]uint8, os.Error) {
	size, err := read4(r);
	if err != nil {
		return []uint8{}, err
	}

	bytes, err := ioutil.ReadAll(io.LimitReader(r, int64(size)));
	if err != nil {
		return []uint8{}, err
	}

	return bytes, nil;
}

func readComplex(r io.Reader) (Term, os.Error) {
	term, err := readTag(r);

	if err != nil {
		return term, err
	}

	switch kind := term.(type) {
	case Atom:
		switch kind {
		case NilAtom:
			return nil, nil
		case TrueAtom:
			return true, nil
		case FalseAtom:
			return false, nil
		}
	}

	return term, nil;
}

func readTag(r io.Reader) (Term, os.Error) {
	tag, err := read1(r);
	if err != nil {
		return nil, err
	}

	switch tag {
	case SmallIntTag:
		return readSmallInt(r)
	case IntTag:
		return readInt(r)
	case SmallBignumTag:
		return nil, ErrUnknownType
	case LargeBignumTag:
		return nil, ErrUnknownType
	case FloatTag:
		return nil, ErrUnknownType
	case AtomTag:
		return readAtom(r)
	case SmallTupleTag:
		return readSmallTuple(r)
	case LargeTupleTag:
		return nil, ErrUnknownType
	case NilTag:
		return readNil(r)
	case StringTag:
		return readString(r)
	case ListTag:
		return readList(r)
	case BinTag:
		return readBin(r)
	}

	return nil, ErrUnknownType;
}

func DecodeFrom(r io.Reader) (Term, os.Error) {
	version, err := read1(r);

	if err != nil {
		return nil, err
	}

	// check protocol version
	if version != VersionTag {
		return nil, ErrBadMagic
	}

	return readTag(r);
}

func Decode(data []byte) (Term, os.Error)	{ return DecodeFrom(bytes.NewBuffer(data)) }
