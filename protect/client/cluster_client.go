package client

import (
	"github.com/benschw/go-protect/protect/api"
	"log"
)

var _ = log.Print

type ClusterClient struct {
	Host string
}

func (c *ClusterClient) GetMembers() (map[string]*api.Member, error) {
	var members map[string]*api.Member

	url := c.Host + "/cluster/member"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return members, err
	}
	err = processResponseEntity(r, &members, 200)
	return members, err
}

func (c *ClusterClient) GetPeers() (map[string]*api.Member, error) {
	var peers map[string]*api.Member

	url := c.Host + "/cluster/peer"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return peers, err
	}
	err = processResponseEntity(r, &peers, 200)
	return peers, err
}

func (c *ClusterClient) GetLeader() (api.Member, error) {
	var leader api.Member

	url := c.Host + "/cluster/leader"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return leader, err
	}
	err = processResponseEntity(r, &leader, 200)
	return leader, err
}
