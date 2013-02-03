build:
	go build

install:
	go install

format:
	go fmt

man: man/gobert.3

man/gobert.3:
	ron -b man/gobert.3.ron
