include $(GOROOT)/src/Make.$(GOARCH)

TARG=bert
GOFILES=\
	decode.go
GOTESTFILES=\
	decode_test.go

include $(GOROOT)/src/Make.pkg
