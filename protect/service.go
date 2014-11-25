package protect

import (
	//	"fmt"
	"github.com/benschw/go-protect/raft/command"
	"github.com/benschw/go-protect/raft/db"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
	"log"
	"os"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

type Service struct {
}

func (s *Service) Run(cfg Config) error {
	leader := ""

	raft.RegisterCommand(&command.WriteCommand{})

	if err := os.MkdirAll("data", 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	db := db.New()

	raftServer := server.New("./data", db, "localhost", 8081)

	runRaftServer(raftServer, leader)

	todoResource := &KeyResource{repo: Repository{db: db, raftServer: raftServer}}

	r := gin.Default()

	// r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	r.Run(cfg.SvcHost)

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
