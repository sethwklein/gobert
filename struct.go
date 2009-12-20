package bert

import (
	"bytes";
	"io";
	"os";
	"reflect";
)

func Unmarshal(data []byte, val interface{}) (err os.Error) {
	result, _ := Decode(data);

	value := reflect.NewValue(val).(*reflect.PtrValue).Elem();

	switch v := value.(type) {
	case *reflect.StructValue:
		slice := reflect.NewValue(result).(*reflect.SliceValue);
		for i := 0; i < slice.Len(); i++ {
			e := slice.Elem(i).(*reflect.InterfaceValue).Elem();
			v.Field(i).SetValue(e);
		}
	}

	return nil;
}

func Marshal(w io.Writer, val interface{}) os.Error {
	buf := bytes.NewBuffer([]byte{});
	write1(buf, VersionTag);
	err := writeTag(buf, reflect.NewValue(val));
	buf.WriteTo(w);
	return err;
}
