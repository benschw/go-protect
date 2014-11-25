package protect

import (
	"github.com/benschw/go-protect/api"
	"github.com/benschw/go-protect/raft/command"
	"github.com/benschw/go-protect/raft/db"
	"github.com/benschw/go-protect/raft/server"
)

type Repository struct {
	db         *db.DB
	raftServer *server.Server
}

func (r *Repository) Set(key api.Key) (api.Key, error) {

	// Execute the command against the Raft server.
	if _, err := r.raftServer.RaftServer().Do(command.NewWriteCommand(key.Id, key.Key)); err != nil {
		return key, err
	}

	return key, nil
}

func (r *Repository) Get(id string) (api.Key, error) {

	keyStr := r.db.Get(id)
	return api.Key{Id: id, Key: keyStr}, nil
}
