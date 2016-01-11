package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type KeyMapper func(string) (tbl string, kcol string, vcol string, kval string)

type DynamoModule struct {
	config    *DynamoCfg
	client    *dynamodb.DynamoDB
	keyMapper KeyMapper
}

type DynamoCfg struct {
	Region     string
	Endpoint   string
	DisableSSL bool
}
