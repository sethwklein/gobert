package bert

import (
	"os";
)

const (
	VersionTag	= 131;
	SmallIntTag	= 97;
	IntTag		= 98;
	SmallBignumTag	= 110;
	LargeBignumTag	= 111;
	FloatTag	= 99;
	AtomTag		= 100;
	SmallTupleTag	= 104;
	LargeTupleTag	= 105;
	NilTag		= 106;
	StringTag	= 107;
	ListTag		= 108;
	BinTag		= 109;
)

type Atom string

type Term interface{}

type Error struct {
	os.ErrorString;
}
