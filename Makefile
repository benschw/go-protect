default: build

clean:
	rm -rf node?
	rm -rf test?

deps:
	go get

build:
	go build

test:
	go test
