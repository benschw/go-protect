package protect

import (
	"github.com/benschw/go-protect/protect/api"
	"github.com/gin-gonic/gin"
	"log"
)

type KeyResource struct {
	repo Repository
}

func (r *KeyResource) CreateKey(c *gin.Context) {
	var key api.Key

	if !c.Bind(&key) {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}

	key, err := r.repo.Set(key)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, api.NewError("problem writing to raft server"))
		return
	}

	c.JSON(201, key)
}

func (r *KeyResource) GetKey(c *gin.Context) {
	idStr := c.Params.ByName("id")

	key, err := r.repo.Get(idStr)
	if err != nil {
		log.Fatal(err)
		c.JSON(404, api.NewError("key not found"))
		return
	}

	c.JSON(200, key)
}
