package main

import (
	"fmt"
	"github.com/benschw/go-protect/protect"
	"github.com/benschw/go-protect/raft/command"
	"github.com/goraft/raft"
	"log"
	"math/rand"
	"time"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func NewMemCluster() protect.Config {
	cfg1 := protect.Config{
		RaftHost:  "localhost",
		RaftPort:  5000,
		ApiHost:   "localhost",
		ApiPort:   6000,
		DataDir:   "test1",
		JoinAddr:  "",
		Bootstrap: true,
	}
	cfg2 := protect.Config{
		RaftHost:  "localhost",
		RaftPort:  5001,
		ApiHost:   "localhost",
		ApiPort:   6001,
		DataDir:   "test2",
		JoinAddr:  fmt.Sprintf("%s:%d", cfg1.RaftHost, cfg1.RaftPort),
		Bootstrap: false,
	}
	cfg3 := protect.Config{
		RaftHost:  "localhost",
		RaftPort:  5002,
		ApiHost:   "localhost",
		ApiPort:   6002,
		DataDir:   "test3",
		JoinAddr:  fmt.Sprintf("%s:%d", cfg1.RaftHost, cfg1.RaftPort),
		Bootstrap: false,
	}

	rand.Seed(time.Now().UnixNano())

	raft.RegisterCommand(&command.WriteCommand{})

	go run(cfg1)
	go run(cfg2)
	go run(cfg3)

	time.Sleep(300 * time.Millisecond)
	return cfg1
}

func run(cfg protect.Config) {
	svc := protect.Service{}
	if err := svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
