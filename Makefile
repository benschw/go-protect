default: build

build:
	go build
	cd cmd/todo; \
	go build

test:
	./go-protect --config config.yaml server & \
	pid=$$!; \
	go test ./...; \
	kill $$pid
