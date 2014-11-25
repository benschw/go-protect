package protect

import (
	//	"fmt"
	"fmt"
	"github.com/benschw/go-protect/raft/command"
	"github.com/benschw/go-protect/raft/db"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
	"log"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	RaftHost string
	RaftPort int
	ApiHost  string
	ApiPort  int
	DataDir  string
	JoinAddr string
}

type Service struct {
}

func (s *Service) Run(cfg Config) error {

	rand.Seed(time.Now().UnixNano())

	raft.RegisterCommand(&command.WriteCommand{})

	if err := os.MkdirAll(cfg.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	db := db.New()

	raftServer := server.New(cfg.DataDir, db, cfg.RaftHost, cfg.RaftPort)

	runRaftServer(raftServer, cfg.JoinAddr)

	todoResource := &KeyResource{repo: Repository{db: db, raftServer: raftServer}}

	r := gin.Default()

	// r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	r.Run(fmt.Sprintf("%s:%d", cfg.ApiHost, cfg.ApiPort))

	return nil
}

func runRaftServer(raftServer *server.Server, leader string) {

	if err := raftServer.Start(); err != nil {
		log.Fatal(err)
	}
	if err := raftServer.Join(leader); err != nil {
		log.Fatal(err)
	}

	go startRaftHttpServer(raftServer)
}

func startRaftHttpServer(raftServer *server.Server) {

	log.Fatal(raftServer.ListenAndServe())

}
