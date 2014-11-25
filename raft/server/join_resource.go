package server

import (
	"encoding/json"
	"github.com/goraft/raft"
	"net/http"
)

type JoinResource struct {
	raftServer raft.Server
}

func (r *JoinResource) joinHandler(w http.ResponseWriter, req *http.Request) {
	command := &raft.DefaultJoinCommand{}

	if err := json.NewDecoder(req.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := r.raftServer.Do(command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
