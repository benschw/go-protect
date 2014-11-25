package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goraft/raft"
	"log"
	"net/http"
)

var _ = log.Print

type RaftMembershipClient struct {
	Host string
}

func (c *RaftMembershipClient) Join(name string, connectionString string) error {
	command := &raft.DefaultJoinCommand{
		Name:             name,
		ConnectionString: connectionString,
	}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(command)
	resp, err := http.Post(fmt.Sprintf("%s/join", c.Host), "application/json", &b)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
