package protect

import (
	// "github.com/benschw/go-protect/api"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
)

type ClusterResource struct {
	config     Config
	raftServer *server.Server
}

func (r *ClusterResource) GetMembers(c *gin.Context) {

	var members map[string]*raft.Peer = r.raftServer.RaftServer().Peers()

	name := r.raftServer.RaftServer().Leader()
	members[name] = &raft.Peer{}

	c.JSON(200, members)
}

func (r *ClusterResource) GetPeers(c *gin.Context) {

	var peers map[string]*raft.Peer = r.raftServer.RaftServer().Peers()

	c.JSON(200, peers)
}

func (r *ClusterResource) GetLeader(c *gin.Context) {

	var leader string = r.raftServer.RaftServer().Leader()

	c.JSON(200, leader)
}
