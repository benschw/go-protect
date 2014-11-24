package client

import (
	"github.com/benschw/go-protect/api"
	"log"
)

var _ = log.Print

type ProtectClient struct {
	Host string
}

func (tc *ProtectClient) CreateKey(id string, keyStr string) (api.Key, error) {
	var respKey api.Key
	key := api.Key{Id: id, Key: keyStr}

	url := tc.Host + "/key"
	r, err := makeRequest("POST", url, key)
	if err != nil {
		return respKey, err
	}
	err = processResponseEntity(r, &respKey, 201)
	return respKey, err
}

func (tc *ProtectClient) GetKey(id string) (api.Key, error) {
	var respKey api.Key

	url := tc.Host + "/key/" + id
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return respKey, err
	}
	err = processResponseEntity(r, &respKey, 200)
	return respKey, err
}
