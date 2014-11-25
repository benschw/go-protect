package service

import (
	//	"fmt"
	"github.com/benschw/go-protect/db"
	//	"github.com/benschw/go-protect/raft/client"
	"github.com/benschw/go-protect/raft/command"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
	"log"
	"os"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/jinzhu/gorm"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

type ProtectService struct {
}

func (s *ProtectService) Run(cfg Config) error {
	leader := ""

	raft.RegisterCommand(&command.WriteCommand{})

	if err := os.MkdirAll("data", 0744); err != nil {
		log.Fatalf("Unable to create path: %v", err)
	}

	raftServer := server.New("./data", "localhost", 8081)

	if err := raftServer.Start(); err != nil {
		log.Fatal(err)
	}
	if err := raftServer.Join(leader); err != nil {
		log.Fatal(err)
	}

	go startRaftServer(raftServer)

	db := db.New()

	todoResource := &KeyResource{db: db}

	r := gin.Default()

	// r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	r.Run(cfg.SvcHost)

	return nil
}

func startRaftServer(raftServer *server.Server) {

	log.Fatal(raftServer.ListenAndServe())

}
