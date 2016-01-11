package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/qzaidi/redamo/store"
)

type DynamoModule struct {
	store.Store
	config    *DynamoCfg
	client    *dynamodb.DynamoDB
	keyMapper store.KeyMapper
}

type DynamoCfg struct {
	Region     string
	Endpoint   string
	DisableSSL bool
}
