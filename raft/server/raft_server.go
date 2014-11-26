package server

import (
	"errors"
	"fmt"
	"github.com/benschw/go-protect/raft/client"
	"github.com/benschw/go-protect/raft/db"
	"github.com/goraft/raft"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// The raftd server is a combination of the Raft server and an HTTP
// server which acts as the transport.
type Server struct {
	name       string
	host       string
	port       int
	path       string
	router     *mux.Router
	raftServer raft.Server
	httpServer *http.Server
	db         *db.DB
	mutex      sync.RWMutex
}

// Creates a new server.
func New(path string, db *db.DB, host string, port int) *Server {
	s := &Server{
		host:   host,
		port:   port,
		path:   path,
		db:     db,
		router: mux.NewRouter(),
	}

	// Read existing name or generate a new one.
	if b, err := ioutil.ReadFile(filepath.Join(path, "name")); err == nil {
		s.name = string(b)
	} else {
		s.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		if err = ioutil.WriteFile(filepath.Join(path, "name"), []byte(s.name), 0644); err != nil {
			panic(err)
		}
	}

	return s
}

// Returns the backing raft.Server
func (s *Server) RaftServer() raft.Server {
	return s.raftServer
}

// Starts the Raft Server
func (s *Server) Start() error {
	var err error

	// Initialize and start Raft server.
	transporter := raft.NewHTTPTransporter("/raft", 200*time.Millisecond)
	s.raftServer, err = raft.NewServer(s.name, s.path, transporter, nil, s.db, "")
	if err != nil {
		return err
	}
	transporter.Install(s.raftServer, s)
	s.raftServer.Start()

	return nil
}

// Joins new cluster as the first node
func (s *Server) Bootstrap() error {
	if !s.raftServer.IsLogEmpty() {
		return errors.New("Cannot bootstrap new cluster with an existing log")
	}

	_, err := s.raftServer.Do(&raft.DefaultJoinCommand{
		Name:             s.raftServer.Name(),
		ConnectionString: s.connectionString(),
	})
	return err
}

// Joins an existing cluster using HTTP client
func (s *Server) Join(leader string) error {
	if s.IsInitialized() {
		return errors.New(fmt.Sprintf("Node already a part of a cluster. Cannot join leader %s", leader))
	}

	// use the raft client to join existing cluster
	raftClient := client.RaftMembershipClient{Host: fmt.Sprintf("http://%s", leader)}
	return raftClient.Join(s.raftServer.Name(), s.connectionString())
}

// Checks is node is already a part of a cluster
func (s *Server) IsInitialized() bool {
	return !s.raftServer.IsLogEmpty()
}

// Starts the http server.
func (s *Server) ListenAndServe() error {

	// Initialize and start HTTP server.
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	joinResource := JoinResource{raftServer: s.raftServer}

	s.router.HandleFunc("/join", joinResource.joinCluster).Methods("POST")

	return s.httpServer.ListenAndServe()
}

// This is a hack around Gorilla mux not providing the correct net/http
// HandleFunc() interface.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}

// Returns the connection string.
func (s *Server) connectionString() string {
	return fmt.Sprintf("http://%s:%d", s.host, s.port)
}
