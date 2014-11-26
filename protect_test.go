package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	"github.com/benschw/go-protect/protect"
	. "gopkg.in/check.v1"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	cfg protect.Config
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.cfg = NewMemCluster()
}

// func (s *MySuite) TestHelloWorld(c *C) {
// 	c.Assert(42, Equals, "42")
// 	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
// 	c.Check(42, Equals, 42)
// }

func (s *MySuite) TestCreateKey(c *C) {

	// given
	client := client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", s.cfg.ApiHost, s.cfg.ApiPort)}

	idStr := "foo"
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"

	// when
	key, err := client.CreateKey(idStr, keyStr)

	//then
	c.Assert(err, Equals, nil)
	c.Assert(key.Id, Equals, idStr)
	c.Assert(key.Key, Equals, keyStr)
}

func (s *MySuite) TestGetKey(c *C) {

	// given
	client := client.ProtectClient{Host: fmt.Sprintf("http://%s:%d", s.cfg.ApiHost, s.cfg.ApiPort)}

	idStr := "foo"
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"
	client.CreateKey(idStr, keyStr)

	// when
	key, err := client.GetKey(idStr)

	// then
	c.Assert(err, Equals, nil)
	c.Assert(key.Id, Equals, idStr)
	c.Assert(key.Key, Equals, keyStr)
}

func (s *MySuite) TestMgmtGetPeers(c *C) {

	// given
	client := client.MgmtClient{Host: fmt.Sprintf("http://%s:%d", s.cfg.ApiHost, s.cfg.ApiPort)}

	// when
	peers, err := client.GetPeers()

	// then
	c.Assert(err, Equals, nil)
	c.Assert(len(peers), Equals, 2)
}

func (s *MySuite) TestMgmtGetLeader(c *C) {

	// given
	client := client.MgmtClient{Host: fmt.Sprintf("http://%s:%d", s.cfg.ApiHost, s.cfg.ApiPort)}

	// when
	leader, err := client.GetLeader()

	// then
	c.Assert(err, Equals, nil)
	c.Assert(leader, Not(Equals), "")
}
