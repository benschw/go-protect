package protect

import (
	//	"fmt"
	"fmt"
	"github.com/benschw/go-protect/raft/db"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type Service struct {
}

// Configure and start the service
func (s *Service) Run(cfg Config) error {
	raftConnectionString := fmt.Sprintf("%s:%d", cfg.RaftHost, cfg.RaftPort)
	apiConnectionString := fmt.Sprintf("%s:%d", cfg.ApiHost, cfg.ApiPort)

	log.Println("=======================================================")
	log.Printf("  Raft Server: http://%s", raftConnectionString)

	log.Printf("  API Server: http://%s", apiConnectionString)
	log.Println()

	// Configuring
	if err := os.MkdirAll(cfg.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	db := db.New()

	raftServer := server.New(cfg.DataDir, db, cfg.RaftHost, cfg.RaftPort)

	todoResource := &KeyResource{repo: Repository{db: db, raftServer: raftServer}}
	mgmtResource := &MgmtResource{raftServer: raftServer}

	// Start API HTTP Server
	log.Println("Initializing API Server")

	r := gin.Default()

	r.GET("/mgmt/peer", mgmtResource.GetPeers)
	r.GET("/mgmt/leader", mgmtResource.GetLeader)

	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	go r.Run(apiConnectionString)

	// Start Raft Server
	log.Printf("Initializing Raft Server: %s", cfg.DataDir)

	if err := raftServer.Start(); err != nil {
		return err
	}

	switch true {
	case cfg.Bootstrap:
		log.Println("Initializing new cluster")

		if err := raftServer.Bootstrap(); err != nil {
			return err
		}
		break
	case cfg.JoinAddr != "" && !raftServer.IsInitialized():
		log.Println("Attempting to join leader:", cfg.JoinAddr)

		if err := raftServer.Join(cfg.JoinAddr); err != nil {
			return err
		}
		break
	case cfg.JoinAddr != "" && raftServer.IsInitialized():
		log.Println("Log already exists, will not attempt join")
		log.Println("Recovering from log")
	default:
		log.Println("Recovering from log")
	}

	log.Fatal(raftServer.ListenAndServe())

	return nil
}
