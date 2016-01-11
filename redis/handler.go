package redis

import (
	redis "github.com/qzaidi/go-redis-server"
  store "github.com/qzaidi/redamo/store"
)

// over-ride the default handler
// TODO: convert store to be an interface
type RedamoHandler struct {
	redis.DefaultHandler
  store *store.DynamoModule
}

func (h *RedamoHandler) Info() ([]byte, error) {
	return []byte(
		`#Server
Version 0.0.1
`), nil
}

func (h *RedamoHandler) Get(key string) ([]byte, error) {
  res,err := h.store.Get(key)
  return res,err
}

func NewRedamoServer(port int,store *store.DynamoModule) (*redis.Server, error) {
	redamo := &RedamoHandler{}
  redamo.store = store
	return redis.NewServer(redis.DefaultConfig().Port(port).Handler(redamo))
}
