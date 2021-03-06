package redis

import (
	"expvar"
	"fmt"
	redis "github.com/qzaidi/go-redis-server"
	store "github.com/qzaidi/redamo/store"

	"time"
)

// over-ride the default handler
type RedamoHandler struct {
	redis.DefaultHandler
	store store.Store
	start time.Time
	tcp   *expvar.Int // total commands processed
}

func (h *RedamoHandler) Info() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`#Server
redamo_version 0.0.1
uptime_in_seconds: %d
#Stats
total_commands_processed: %s
`, int(time.Since(h.start).Seconds()), h.tcp.String())), nil
}

func (h *RedamoHandler) Get(key string) ([]byte, error) {
	res, err := h.store.Get(key)
	h.tcp.Add(1)
	return res, err
}

func (h *RedamoHandler) Set(key string, val []byte) error {
	h.tcp.Add(1)
	return h.store.Set(key, val)
}

func (h *RedamoHandler) Incrby(key string, incr []byte) (int, error) {
	h.tcp.Add(1)
	return h.store.Incrby(key, incr)
}

func (h *RedamoHandler) Incr(key string) (int, error) {
	h.tcp.Add(1)
	return h.store.Incr(key)
}

func NewRedamoServer(port int, store store.Store) (*redis.Server, error) {
	redamo := &RedamoHandler{}
	redamo.store = store
	redamo.start = time.Now()
	redamo.tcp = expvar.NewInt("tcp")
	return redis.NewServer(redis.DefaultConfig().Port(port).Handler(redamo))
}
