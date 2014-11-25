	./go-protect -raft localhost:5000 -api localhost:6000 -data node1 serve	
	./go-protect -raft localhost:5001 -api localhost:6001 -data node2 -join localhost:5000 serve
	./go-protect -raft localhost:5002 -api localhost:6002 -data node3 -join localhost:5000 serve

	curl -H "Content-Type: application/json" -d '{"id":"xyz","key":"blahblahblah"}' http://localhost:6000/key
	curl http://localhost:6000/key/xyz

## REST Microservice in Go for Gin

https://github.com/goraft/raftd

Example (seed) project for microservices in go using the [Gin](http://gin-gonic.github.io/gin/) web framework.


See the [blog post](http://txt.fliglio.com) for a walk through.

### Hacking

#### Build Service
	
	make build

#### Build the Database

	mysql -u root -p -e 'Create Database Todo;'

	./cmd/server/server --config config.yaml migratedb

There's also a `make` target: `make migrate`, but you still need to create the database by hand.

#### Run the Service

	./cmd/server/server --config config.yaml server

#### Testing
The tests leverage a running instance of the server. This is automated with the `Make` target

	make test
