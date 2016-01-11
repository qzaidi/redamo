package store

type KeyMapper func(string) (tbl string, kcol string, vcol string, kval string)

type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
}
