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

	todoResource := &KeyResource{repo: Repository{db: db, raftServer: raftServer}}
	mgmtResource := &MgmtResource{raftServer: raftServer}

	// Start API HTTP Server
	r := gin.Default()

	r.GET("/mgmt/peer", mgmtResource.GetPeers)
	r.GET("/mgmt/leader", mgmtResource.GetLeader)

	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	go r.Run(fmt.Sprintf("%s:%d", cfg.ApiHost, cfg.ApiPort))

	// Start Raft Server
	if err := raftServer.Start(); err != nil {
		log.Fatal(err)
	}
	if err := raftServer.Join(cfg.JoinAddr); err != nil {
		log.Fatal(err)
	}

	log.Fatal(raftServer.ListenAndServe())

	return nil
}
