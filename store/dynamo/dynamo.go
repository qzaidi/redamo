package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/qzaidi/redamo/store"
	logging "gopkg.in/tokopedia/logging.v1"
	"log"
)

// here, support a function that maps a given redis key to a dynamodb table key and value field
func NewDynamoModule(keymap store.KeyMapper) *DynamoModule {

	module := &DynamoModule{}
	cfgdir := "etc"
	cfg := module.config
	ok := logging.ReadModuleConfig(cfg, cfgdir, "dynamo") || logging.ReadModuleConfig(cfg, ".", "dynamo")

	sess := session.New(&aws.Config{Region: aws.String("ap-southeast-1")})

	if !ok {
		log.Println("failed to read dynamo config, using defaults")
	} else {
		sess = session.New(&aws.Config{
			Region:     aws.String(cfg.Region),
			Endpoint:   aws.String(cfg.Endpoint),
			DisableSSL: aws.Bool(cfg.DisableSSL),
		})
	}

	module.client = dynamodb.New(sess)
	module.keyMapper = keymap

	return module
}

// initially, we only support string type for key and value, to keep things simple
func (d *DynamoModule) Set(key string, val []byte) error {

	// get the table to be used, key column name,value column name, and the key value

	tbl, kcol, vcol, kval := d.keyMapper(key)

	params := &dynamodb.UpdateItemInput{
		Key:              map[string]*dynamodb.AttributeValue{},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{},
		TableName:        aws.String(tbl),
	}

	params.Key[kcol] = &dynamodb.AttributeValue{S: aws.String(kval)}

	params.AttributeUpdates[vcol] = &dynamodb.AttributeValueUpdate{
		Action: aws.String("PUT"),
		Value: &dynamodb.AttributeValue{
			S: aws.String(string(val)),
		},
	}
	_, err := d.client.UpdateItem(params)
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoModule) Get(key string) ([]byte, error) {
	tbl, kcol, vcol, kval := d.keyMapper(key)
	params := &dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{},
		TableName: aws.String(tbl),
	}
	params.Key[kcol] = &dynamodb.AttributeValue{S: aws.String(kval)}

	resp, err := d.client.GetItem(params)

	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, fmt.Errorf("key attribute not found: %s", kval)
	}

	if resp.Item[vcol] != nil {
		if resp.Item[vcol].S != nil {
			return []byte(*resp.Item[vcol].S), nil
		} else if resp.Item[vcol].N != nil {
			return []byte(*resp.Item[vcol].N), nil
		}
	}

	return nil, fmt.Errorf("value attribute not found: %s", vcol)
}
