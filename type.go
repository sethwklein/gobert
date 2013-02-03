package bert

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

const (
	BertAtom	= Atom("bert");
	NilAtom		= Atom("nil");
	TrueAtom	= Atom("true");
	FalseAtom	= Atom("false");
)

type Term interface{}

type Request struct {
	Kind		Atom;
	Module		Atom;
	Function	Atom;
	Arguments	[]Term;
}
