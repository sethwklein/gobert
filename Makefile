include $(GOROOT)/src/Make.$(GOARCH)

TARG=bert
GOFILES=\
	decode.go\
	struct.go
GOTESTFILES=\
	decode_test.go\
	struct_test.go

include $(GOROOT)/src/Make.pkg

format:
	echo $(GOFILES) $(GOTESTFILES) | xargs gofmt -w
