package main

import (
	"fmt"
	"github.com/benschw/go-protect/protect"
	"github.com/benschw/go-protect/raft/command"
	"github.com/goraft/raft"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

type RandomApiFactory struct{}

func (a *RandomApiFactory) Config(bootstrap bool, leaderConfig *protect.Config) *protect.Config {
	joinAddr := ""
	if leaderConfig != nil {
		joinAddr = fmt.Sprintf("%s:%d", leaderConfig.RaftHost, leaderConfig.RaftPort)
	}

	return &protect.Config{
		RaftHost:  "localhost",
		RaftPort:  a.getRandomPort(),
		ApiHost:   "localhost",
		ApiPort:   a.getRandomPort(),
		DataDir:   a.getRandomDataDir(),
		JoinAddr:  joinAddr,
		Bootstrap: bootstrap,
	}

}
func (a *RandomApiFactory) getRandomPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	addrParts := strings.Split(l.Addr().String(), ":")
	port, _ := strconv.Atoi(addrParts[len(addrParts)-1])

	return port
}
func (a *RandomApiFactory) getRandomDataDir() string {
	nodeStr := fmt.Sprintf("%07x", rand.Int())[0:7]

	return fmt.Sprintf("test-data/%s", nodeStr)
}

func NewMemCluster() protect.Config {
	rand.Seed(time.Now().UnixNano())

	factory := new(RandomApiFactory)
	cfg1 := factory.Config(true, nil)
	cfg2 := factory.Config(false, cfg1)
	cfg3 := factory.Config(false, cfg1)

	raft.RegisterCommand(&command.WriteCommand{})

	go run(*cfg1)
	go run(*cfg2)
	go run(*cfg3)

	time.Sleep(300 * time.Millisecond)
	return *cfg1
}

func run(cfg protect.Config) {
	svc := protect.Service{}
	if err := svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
