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

// Get a random available port
func GetRandomPort() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	addrParts := strings.Split(l.Addr().String(), ":")
	port, _ := strconv.Atoi(addrParts[len(addrParts)-1])

	return port
}

// Get a random available directory
func GetRandomDataDir() string {
	nodeStr := fmt.Sprintf("%07x", rand.Int())[0:7]

	return fmt.Sprintf("test-data/%s", nodeStr)
}

type RandomApiFactory struct{}

// Generate a random config
func (a *RandomApiFactory) Config(bootstrap bool, leaderConfig *protect.Config) *protect.Config {
	joinAddr := ""
	if leaderConfig != nil {
		joinAddr = fmt.Sprintf("%s:%d", leaderConfig.RaftHost, leaderConfig.RaftPort)
	}

	return &protect.Config{
		RaftHost:  "localhost",
		RaftPort:  GetRandomPort(),
		ApiHost:   "localhost",
		ApiPort:   GetRandomPort(),
		DataDir:   GetRandomDataDir(),
		JoinAddr:  joinAddr,
		Bootstrap: bootstrap,
	}

}

type TestCluster struct {
	leaderConfig    protect.Config
	followerConfigs []protect.Config
}

// Build a random TestCluster and return configs used to create
func NewTestCluster(nodes int) *TestCluster {
	rand.Seed(time.Now().UnixNano())
	raft.RegisterCommand(&command.WriteCommand{})

	factory := new(RandomApiFactory)

	leaderConfig := factory.Config(true, nil)
	go runTestServer(*leaderConfig)

	followers := make([]protect.Config, nodes-1, nodes-1)

	for i := 0; i < nodes-1; i++ {
		cfg := factory.Config(false, leaderConfig)
		go runTestServer(*cfg)

		followers[i] = *cfg
	}

	time.Sleep(300 * time.Millisecond)
	return &TestCluster{
		leaderConfig:    *leaderConfig,
		followerConfigs: followers,
	}
}

func runTestServer(cfg protect.Config) {
	svc := protect.Service{}
	if err := svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
