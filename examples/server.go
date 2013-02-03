package main

import (
	"log"
	"github.com/josh/gobert"
	"net"
)

func fib(n int) int {
	if n == 1 {
		return 0
	} else if n == 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func handle(c net.Conn) {
	request, _ := bert.UnmarshalRequest(c)

	var response []bert.Term

	if request.Function == bert.Atom("fib") {
		result := fib(request.Arguments[0].(int))
		response = []bert.Term{bert.Atom("reply"), result}
	} else {
		msg := "function '" + request.Function + "' not found"
		error := []bert.Term{bert.Atom("server"), 2, msg, []bert.Term{}}
		response = []bert.Term{bert.Atom("error"), error}
	}

	bert.MarshalResponse(c, response)
	c.Close()
	log.Println("Request handled.")
}

func main() {
	address := ":8000"
	l, _ := net.Listen("tcp", address)
	log.Printf("Listening for connections at %v...\n", address)

	for {
		c, _ := l.Accept()
		log.Println("Dispatching handler...")
		go handle(c)
	}

	l.Close()
}
