package bert

import (
	"bytes";
	"io";
	"os";
	"reflect";
)

func UnmarshalFrom(r io.Reader, val interface{}) (err os.Error) {
	result, _ := DecodeFrom(r);

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

func Unmarshal(data []byte, val interface{}) (err os.Error) {
	return UnmarshalFrom(bytes.NewBuffer(data), val)
}

func UnmarshalRequest(r io.Reader) (Request, os.Error) {
	var req Request;

	size, err := read4(r);
	if err != nil {
		return req, err
	}

	err = UnmarshalFrom(io.LimitReader(r, int64(size)), &req);

	return req, err;
}

func Marshal(w io.Writer, val interface{}) os.Error {
	return EncodeTo(w, val)
}

func MarshalResponse(w io.Writer, val interface{}) (err os.Error) {
	resp, err := Encode(val);

	write4(w, uint32(len(resp)));
	w.Write(resp);

	return;
}
