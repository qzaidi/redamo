package store

type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
}
