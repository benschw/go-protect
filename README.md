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
