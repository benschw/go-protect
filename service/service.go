package service

import (
	"github.com/benschw/go-protect/db"
	"github.com/gin-gonic/gin"
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
	db := db.New()

	todoResource := &KeyResource{db: db}

	raftServer := raft.New("./foo", "localhost", "8080")

	r := gin.Default()

	// r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/key/:id", todoResource.GetKey)
	r.POST("/key", todoResource.CreateKey)

	r.Run(cfg.SvcHost)

	return nil
}
