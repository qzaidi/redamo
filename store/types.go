package store

type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Incrby(string, []byte) (int, error)
	Incr(string) (int, error)
}
