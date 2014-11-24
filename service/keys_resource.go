package service

import (
	"github.com/benschw/go-protect/api"
	"github.com/benschw/go-protect/db"
	"github.com/gin-gonic/gin"
)

type KeyResource struct {
	db *db.DB
}

func (tr *KeyResource) CreateKey(c *gin.Context) {
	var key api.Key

	if !c.Bind(&key) {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}

	tr.db.Put(key.Id, key.Key)

	c.JSON(201, key)
}

// func (tr *KeyResource) GetAllKeys(c *gin.Context) {
// 	var todos []api.Todo

// 	tr.db.Order("created desc").Find(&todos)

// 	c.JSON(200, todos)
// }

func (tr *KeyResource) GetKey(c *gin.Context) {
	idStr := c.Params.ByName("id")

	keyStr := tr.db.Get(idStr)
	key := api.Key{Id: idStr, Key: keyStr}

	c.JSON(200, key)

	//	c.JSON(404, gin.H{"error": "not found"})
}
