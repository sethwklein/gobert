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

func UnmarshalRequest(data []byte) (Request, os.Error) {
	var req Request;

	buf := bytes.NewBuffer(data);

	size, err := read4(buf);
	if err != nil {
		return req, err
	}

	err = Unmarshal(buf.Bytes()[0:size], &req);

	return req, err;
}

func Marshal(w io.Writer, val interface{}) os.Error {
	write1(w, VersionTag);
	err := writeTag(w, reflect.NewValue(val));
	return err;
}

func MarshalResponse(w io.Writer, val interface{}) (err os.Error) {
	resp, err := Encode(val);

	write4(w, uint32(len(resp)));
	w.Write(resp);

	return;
}
