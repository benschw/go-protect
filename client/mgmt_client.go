package client

import (
	"github.com/goraft/raft"
	"log"
)

var _ = log.Print

type MgmtClient struct {
	Host string
}

func (c *MgmtClient) GetPeers() (map[string]*raft.Peer, error) {
	var peers map[string]*raft.Peer

	url := c.Host + "/mgmt/peer"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return peers, err
	}
	err = processResponseEntity(r, &peers, 200)
	return peers, err
}

func (c *MgmtClient) GetLeader() (string, error) {

	url := c.Host + "/mgmt/leader"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	leader, err := processResponseBytes(r, 200)
	return string(leader[:]), err
}
