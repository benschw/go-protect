package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestCreateKey(t *testing.T) {

	// given
	client := client.ProtectClient{Host: "http://localhost:6000"}

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
	client := client.ProtectClient{Host: "http://localhost:6000"}
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
	client := client.MgmtClient{Host: "http://localhost:6000"}

	// when
	peers, err := client.GetPeers()

	// then
	if err != nil {
		t.Error(err)
	}
	if len(peers) != 2 {
		t.Error("wrong number of peers")
	}
}

func TestMgmtGetLeader(t *testing.T) {

	// given
	client := client.MgmtClient{Host: "http://localhost:6000"}

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
