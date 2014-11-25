default: build

build:
	go build
	cd cmd/todo; \
	go build

test:
	./go-protect serve & \
	pid=$$!; \
	go test ./...; \
	kill $$pid
