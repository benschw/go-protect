package protect

import (
	"fmt"
	"github.com/benschw/go-protect/protect/api"
	"github.com/benschw/go-protect/raft/server"
	"github.com/gin-gonic/gin"
	"github.com/goraft/raft"
)

type ClusterResource struct {
	config     Config
	raftServer *server.Server
}

func (r *ClusterResource) GetMembers(c *gin.Context) {

	c.JSON(200, r.getAllMembers())
}

func (r *ClusterResource) GetPeers(c *gin.Context) {

	peers := r.raftServer.RaftServer().Peers()

	members := r.toMembers(peers)

	c.JSON(200, members)
}

func (r *ClusterResource) GetLeader(c *gin.Context) {

	var leader string = r.raftServer.RaftServer().Leader()

	members := r.getAllMembers()

	c.JSON(200, members[leader])
}

func (r *ClusterResource) getAllMembers() map[string]*api.Member {
	var peers map[string]*raft.Peer = r.raftServer.RaftServer().Peers()

	name := r.raftServer.RaftServer().Name()
	peers[name] = &raft.Peer{
		Name:             name,
		ConnectionString: fmt.Sprintf("http://%s:%d", r.config.RaftHost, r.config.RaftPort),
	}

	return r.toMembers(peers)
}

func (r *ClusterResource) toMembers(peers map[string]*raft.Peer) map[string]*api.Member {
	var members = make(map[string]*api.Member, len(peers))

	for _, member := range peers {
		members[member.Name] = r.toMember(member)
	}
	return members
}
func (r *ClusterResource) toMember(member *raft.Peer) *api.Member {
	return &api.Member{
		Name:                 member.Name,
		RaftConnectionString: member.ConnectionString,
		ApiConnectionString:  "",
	}
}
