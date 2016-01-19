package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/koding/cache"
	"github.com/qzaidi/redamo/store"
)

type DynamoModule struct {
	store.Store
	config    DynamoCfg
	client    *dynamodb.DynamoDB
	keyMapper KeyMapper
	cache     *cache.MemoryTTL
}

// a mapper maps a redis prefix to dynamo table and key
type Mapper struct {
	Table  string
	Kcol   string
	Vcol   string
	Ktype  string
	Vtype  string
	Keyval string
}

type KeyMapper func(string) *Mapper

type DynamoCfg struct {
	Server struct {
		Region        string
		Endpoint      string
		DisableSSL    bool
		CacheDuration int
	}
	Keymap map[string]*Mapper
}
