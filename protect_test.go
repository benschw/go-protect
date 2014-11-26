package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	. "gopkg.in/check.v1"
	"log"
	"sort"
	"testing"
	"time"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	members int
	cluster *TestCluster

	protectClient  client.ProtectClient
	fProtectClient client.ProtectClient
	clusterClient  client.ClusterClient
	fClusterClient client.ClusterClient
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.members = 3
	s.cluster = NewTestCluster(s.members)

	s.protectClient = client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", s.cluster.leaderConfig.ApiHost, s.cluster.leaderConfig.ApiPort)}
	s.fProtectClient = client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", s.cluster.followerConfigs[0].ApiHost, s.cluster.followerConfigs[0].ApiPort)}

	s.clusterClient = client.ClusterClient{Host: fmt.Sprintf("http://%s:%d", s.cluster.leaderConfig.ApiHost, s.cluster.leaderConfig.ApiPort)}
	s.fClusterClient = client.ClusterClient{Host: fmt.Sprintf("http://%s:%d", s.cluster.followerConfigs[0].ApiHost, s.cluster.followerConfigs[0].ApiPort)}
}

func (s *MySuite) TestProtect_CreateKey(c *C) {

	// given
	idStr := "foo"
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"

	// when
	key, err := s.protectClient.CreateKey(idStr, keyStr)

	//then
	c.Assert(err, Equals, nil)
	c.Assert(key.Id, Equals, idStr)
	c.Assert(key.Key, Equals, keyStr)
}

func (s *MySuite) TestProtect_GetKey(c *C) {

	// given
	idStr := "foo"
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"
	s.protectClient.CreateKey(idStr, keyStr)

	// when
	key, err := s.protectClient.GetKey(idStr)

	// then
	c.Assert(err, Equals, nil)
	c.Assert(key.Id, Equals, idStr)
	c.Assert(key.Key, Equals, keyStr)
}

func (s *MySuite) TestProtect_GetKeyFromFollower(c *C) {

	// given
	idStr := "foo"
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"
	s.protectClient.CreateKey(idStr, keyStr)
	time.Sleep(100 * time.Millisecond) // give it time to become consistent

	// when
	key, err := s.fProtectClient.GetKey(idStr)

	// then
	c.Assert(err, Equals, nil)
	c.Assert(key.Id, Equals, idStr)
	c.Assert(key.Key, Equals, keyStr)
}

func (s *MySuite) TestCluster_GetMembers(c *C) {
	expected := make([]string, s.members, s.members)
	for i, cfg := range s.cluster.followerConfigs {
		expected[i] = fmt.Sprintf("http://%s:%d", cfg.RaftHost, cfg.RaftPort)
	}
	expected[s.members-1] = fmt.Sprintf("http://%s:%d", s.cluster.leaderConfig.RaftHost, s.cluster.leaderConfig.RaftPort)
	sort.Strings(expected)

	// when
	members, err := s.clusterClient.GetMembers()

	// then
	c.Assert(err, Equals, nil)
	c.Assert(len(members), Equals, s.members)

	found := make([]string, len(members), len(members))
	i := 0
	for _, peer := range members {
		found[i] = peer.ConnectionString
		i++
	}
	sort.Strings(found)
	c.Assert(found, DeepEquals, expected)
}

func (s *MySuite) TestCluster_GetPeers(c *C) {

	// when
	peers, err := s.clusterClient.GetPeers()

	// then
	c.Assert(err, Equals, nil)
	c.Assert(len(peers), Equals, s.members-1)
}

func (s *MySuite) TestCluster_GetLeader(c *C) {

	// when
	leader, err := s.clusterClient.GetLeader()

	// then
	c.Assert(err, Equals, nil)
	c.Assert(leader, Not(Equals), "")
}
