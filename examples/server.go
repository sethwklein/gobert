package main

import (
	"bert";
	"net";
)

func fib(n int) int {
	if n == 1 {
		return 0
	} else if n == 2 {
		return 1
	}
	return fib(n-1) + fib(n-2);
}

func handle(c net.Conn) {
	buf := make([]byte, 1024);
	c.Read(buf);
	request, _ := bert.UnmarshalRequest(buf);

	var response []bert.Term;

	if request.Function == bert.Atom("fib") {
		result := fib(request.Arguments[0].(int));
		response = []bert.Term{bert.Atom("reply"), result};
	} else {
		msg := "function '" + request.Function + "' not found";
		error := []bert.Term{bert.Atom("server"), 2, msg, []bert.Term{}};
		response = []bert.Term{bert.Atom("error"), error};
	}

	bert.MarshalResponse(c, response);
	c.Close();
}

func main() {
	l, _ := net.Listen("tcp", ":8000");

	for {
		c, _ := l.Accept();
		go handle(c);
	}

	l.Close();
}
