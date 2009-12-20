include $(GOROOT)/src/Make.$(GOARCH)

TARG=bert
GOFILES=\
	decode.go\
	encode.go\
	struct.go\
	type.go
GOTESTFILES=\
	decode_test.go\
	encode_test.go\
	struct_test.go

include $(GOROOT)/src/Make.pkg

format:
	echo $(GOFILES) $(GOTESTFILES) | xargs gofmt -w
