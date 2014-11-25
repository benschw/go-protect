default: build

clean:
	rm -rf node?

deps:
	go get;
	go get code.google.com/p/go.tools/cmd/cover;
	go get github.com/mattn/goveralls;

build:
	go build

test-all: test test-recover

test:
	./go-protect -raft localhost:5000 -api localhost:6000 -data node1 -bootstrap serve	& \
	pid1=$$!; \
	sleep 1; \
	./go-protect -raft localhost:5001 -api localhost:6001 -data node2 -join localhost:5000 serve & \
	pid2=$$!; \
	sleep 1; \
	./go-protect -raft localhost:5002 -api localhost:6002 -data node3 -join localhost:5000 serve & \
	pid3=$$!; \
	go test -v; \
	kill $$pid1; \
	kill $$pid2; \
	kill $$pid3; \


test-recover:
	./go-protect -raft localhost:5000 -api localhost:6000 -data node1 serve	& \
	pid1=$$!; \
	sleep 1; \
	./go-protect -raft localhost:5001 -api localhost:6001 -data node2 serve & \
	pid2=$$!; \
	sleep 1; \
	./go-protect -raft localhost:5002 -api localhost:6002 -data node3 -join localhost:5000 serve & \
	pid3=$$!; \
	go test -v ; \
	kill $$pid1; \
	kill $$pid2; \
	kill $$pid3; \
