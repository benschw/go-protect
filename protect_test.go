package main

import (
	"fmt"
	"github.com/benschw/go-protect/client"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestCreateKey(t *testing.T) {

	// given
	client := client.ProtectClient{Host: "http://localhost:6000"}

	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"

	// when
	key, err := client.CreateKey("foo", keyStr)

	//then
	if err != nil {
		t.Error(err)
	}

	if key.Id != "foo" && key.Key != keyStr {
		t.Error("returned key not right")
	}

	// cleanup
	//	_ = client.DeleteTodo(todo.Id)
}

func TestGetKey(t *testing.T) {

	// given
	client := client.ProtectClient{Host: "http://localhost:6000"}
	keyStr := "1g34jh142jhg1234j412uyg142iuy124guy142g"
	keyCreate, _ := client.CreateKey("foo", keyStr)
	id := keyCreate.Id

	// when
	key, err := client.GetKey(id)

	// then
	if err != nil {
		t.Error(err)
	}

	if key.Id != "foo" && key.Key != keyStr {
		t.Error("returned todo not right")
	}

	// cleanup
	// _ = client.DeleteTodo(todo.Id)
}
