include $(GOROOT)/src/Make.$(GOARCH)

all: server

server: ../_obj server.6
	$(LD) -L../_obj -o $@ server.6

server.6:
	$(GC) -I../_obj -o $@ server.go

../_obj:
	cd .. && $(MAKE)

clean:
	rm *.[568vqo] server
