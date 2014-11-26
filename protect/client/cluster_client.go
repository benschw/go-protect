package client

import (
	"github.com/goraft/raft"
	"log"
)

var _ = log.Print

type ClusterClient struct {
	Host string
}

func (c *ClusterClient) GetMembers() (map[string]*raft.Peer, error) {
	var members map[string]*raft.Peer

	url := c.Host + "/cluster/member"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return members, err
	}
	err = processResponseEntity(r, &members, 200)
	return members, err
}

func (c *ClusterClient) GetPeers() (map[string]*raft.Peer, error) {
	var peers map[string]*raft.Peer

	url := c.Host + "/cluster/peer"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return peers, err
	}
	err = processResponseEntity(r, &peers, 200)
	return peers, err
}

func (c *ClusterClient) GetLeader() (string, error) {

	url := c.Host + "/cluster/leader"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	leader, err := processResponseBytes(r, 200)
	return string(leader[:]), err
}
