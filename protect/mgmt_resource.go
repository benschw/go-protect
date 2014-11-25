package protect

import (
	// "github.com/benschw/go-protect/api"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
)

type MgmtResource struct {
	raftServer *server.Server
}

func (r *MgmtResource) GetPeers(c *gin.Context) {

	var peers map[string]*raft.Peer = r.raftServer.RaftServer().Peers()

	c.JSON(200, peers)
}

func (r *MgmtResource) GetLeader(c *gin.Context) {

	var leader string = r.raftServer.RaftServer().Leader()

	c.JSON(200, leader)
}
