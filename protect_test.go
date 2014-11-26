package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	"github.com/benschw/go-protect/protect"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

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
