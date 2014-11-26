[![Build Status](https://travis-ci.org/benschw/go-protect.svg)](https://travis-ci.org/benschw/go-protect)
[![GoDoc](http://godoc.org/github.com/benschw/go-protect?status.png)](http://godoc.org/github.com/benschw/go-protect)

## build
	make deps
	make build

## test
	make test-all

## try it out
	# start up a 3 node cluster

	./go-protect -raft localhost:5000 -api localhost:6000 -data data/node1 -bootstrap serve
	./go-protect -raft localhost:5001 -api localhost:6001 -data data/node2 -join localhost:5000 serve
	./go-protect -raft localhost:5002 -api localhost:6002 -data data/node3 -join localhost:5000 serve


	# add/retreive a key

	curl -H "Content-Type: application/json" -d '{"id":"xyz","key":"blahblahblah"}' http://localhost:6000/key
	curl http://localhost:6000/key/xyz


	# kill existing cluster and start them back up again

	killall go-protect
	./go-protect -raft localhost:5000 -api localhost:6000 -data data/node1 serve
	./go-protect -raft localhost:5001 -api localhost:6001 -data data/node2 serve
	./go-protect -raft localhost:5002 -api localhost:6002 -data data/node3 serve


## Links
- [Gin](http://gin-gonic.github.io/gin/)
- [goraft](https://github.com/goraft/raftd)
