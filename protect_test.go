package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	"github.com/benschw/go-protect/protect"
	"github.com/benschw/go-protect/raft/command"
	"github.com/goraft/raft"
	"log"
	"math/rand"
	"testing"
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
	time.Sleep(1000 * time.Millisecond)
	go run(cfg2)
	go run(cfg3)
	time.Sleep(1000 * time.Millisecond)

	return cfg1
}

func run(cfg protect.Config) {
	svc := protect.Service{}
	if err := svc.Run(cfg); err != nil {
		log.Fatal(err)
	}
}

var cfg protect.Config

func init() {
	cfg = NewMemCluster()
}
func TestCreateKey(t *testing.T) {

	// given
	client := client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", cfg.ApiHost, cfg.ApiPort)}

	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"

	// when
	key, err := client.CreateKey("foo", keyStr)

	//then
	if err != nil {
		t.Error(err)
	}

	if key.Id != "foo" && key.Key != keyStr {
		t.Error("returned key not right")
	}
}

func TestGetKey(t *testing.T) {

	// given
	client := client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", cfg.ApiHost, cfg.ApiPort)}
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"
	keyCreate, _ := client.CreateKey("foo", keyStr)
	id := keyCreate.Id

	// when
	key, err := client.GetKey(id)

	// then
	if err != nil {
		t.Error(err)
	}

	if key.Id != "foo" && key.Key != keyStr {
		t.Error("returned todo not right")
	}
}

func TestMgmtGetPeers(t *testing.T) {

	// given
	client := client.MgmtClient{Host: fmt.Sprintf("http://%s:%d", cfg.ApiHost, cfg.ApiPort)}

	// when
	peers, err := client.GetPeers()

	// then
	if err != nil {
		t.Error(err)
	}
	if len(peers) != 2 {
		t.Errorf("wrong number of peers; found: %d", len(peers))
	}
}

func TestMgmtGetLeader(t *testing.T) {

	// given
	client := client.MgmtClient{Host: fmt.Sprintf("http://%s:%d", cfg.ApiHost, cfg.ApiPort)}

	// when
	leader, err := client.GetLeader()

	// then
	if err != nil {
		t.Error(err)
	}
	if leader == "" {
		t.Error("leader name not returned")
	}
}
