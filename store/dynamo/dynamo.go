package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/koding/cache"
	logging "gopkg.in/tokopedia/logging.v1"
	"log"
	"strconv"
	"time"
)

// here, support a function that maps a given redis key to a dynamodb table key and value field
func NewDynamoModule(keymap KeyMapper) *DynamoModule {

	module := &DynamoModule{}
	cfgdir := "/etc"
	cfg := &module.config
	ok := logging.ReadModuleConfig(cfg, cfgdir, "dynamo") || logging.ReadModuleConfig(cfg, ".", "dynamo")

	sess := session.New(&aws.Config{Region: aws.String("ap-southeast-1")})

	if !ok {
		log.Println("failed to read dynamo config, using defaults")
	} else {
		sess = session.New(&aws.Config{
			Region:     aws.String(cfg.Server.Region),
			Endpoint:   aws.String(cfg.Server.Endpoint),
			DisableSSL: aws.Bool(cfg.Server.DisableSSL),
		})
	}

	module.client = dynamodb.New(sess)
	if keymap != nil {
		module.keyMapper = keymap
	} else {
		module.keyMapper = module.defaultMapper
	}

	if cfg.Server.CacheDuration > 0 {
		logging.Debug.Println("activiating cache, TTL", cfg.Server.CacheDuration)
		module.cache = cache.NewMemoryWithTTL(time.Duration(cfg.Server.CacheDuration) * time.Second)
	}

	return module
}

// initially, we only support string type for key and value, to keep things simple
func (d *DynamoModule) Set(key string, val []byte) error {

	// get the table to be used, key column name,value column name, and the key value
	kmap := d.keyMapper(key)
	if kmap == nil {
		return fmt.Errorf("bad key")
	}

	params := &dynamodb.UpdateItemInput{
		Key:              map[string]*dynamodb.AttributeValue{},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{},
		TableName:        aws.String(kmap.Table),
	}

	params.Key[kmap.Kcol] = &dynamodb.AttributeValue{S: aws.String(kmap.Keyval)}

	params.AttributeUpdates[kmap.Vcol] = &dynamodb.AttributeValueUpdate{
		Action: aws.String("PUT"),
		Value:  &dynamodb.AttributeValue{},
	}

	switch kmap.Vtype[0] {
	case 'S':
		params.AttributeUpdates[kmap.Vcol].Value.S = aws.String(string(val))
	case 'N':
		params.AttributeUpdates[kmap.Vcol].Value.N = aws.String(string(val))
	}

	_, err := d.client.UpdateItem(params)
	if err != nil {
		return err
	}

	if d.cache != nil {
		d.cache.Set(key, val)
	}

	return nil
}

func (d *DynamoModule) Get(key string) ([]byte, error) {

	if d.cache != nil {
		val, err := d.cache.Get(key)
		if err == nil && val != nil {
      var res []byte
      ok := false
      if res,ok = val.([]byte); !ok {
        err = fmt.Errorf("Bad value in cache for %s",key)
      }
			return res,err
		}
	}

	kmap := d.keyMapper(key)

	if kmap == nil {
		return nil, fmt.Errorf("bad key")
	}

	params := &dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{},
		TableName: aws.String(kmap.Table),
	}
	params.Key[kmap.Kcol] = &dynamodb.AttributeValue{S: aws.String(kmap.Keyval)}

	resp, err := d.client.GetItem(params)

	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	vcol := kmap.Vcol

	if resp.Item[vcol] != nil {
		var val []byte
		switch kmap.Vtype[0] {
		case 'S':
			if resp.Item[vcol].S == nil {
				return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
			val = []byte(*resp.Item[vcol].S)
		case 'N':
			if resp.Item[vcol].N == nil {
				return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
			val = []byte(*resp.Item[vcol].N)
		}
		if val != nil {
      logging.Debug.Println("saving to cache",key)
			if d.cache != nil {
				d.cache.Set(key, val)
			}
			return val, nil
		}
	}

	return nil, fmt.Errorf("value attribute not found: %s", vcol)
}

func (d *DynamoModule) Incrby(key string, val []byte) (int, error) {
	kmap := d.keyMapper(key)

	if kmap == nil || kmap.Vtype[0] != 'N' {
		return -1, fmt.Errorf("unknown key or incorrect config", key, kmap)
	}

	params := &dynamodb.UpdateItemInput{
		Key:              map[string]*dynamodb.AttributeValue{},
		AttributeUpdates: map[string]*dynamodb.AttributeValueUpdate{},
		TableName:        aws.String(kmap.Table),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	params.Key[kmap.Kcol] = &dynamodb.AttributeValue{S: aws.String(kmap.Keyval)}

	params.AttributeUpdates[kmap.Vcol] = &dynamodb.AttributeValueUpdate{
		Action: aws.String("ADD"),
		Value: &dynamodb.AttributeValue{
			N: aws.String(string(val)),
		},
	}

	data, err := d.client.UpdateItem(params)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	newval, err := strconv.ParseInt(*data.Attributes[kmap.Vcol].N, 10, 64)
	if err != nil {
		return -1, err
	}

	if d.cache != nil {
    // cache always stores []byte
		d.cache.Set(key, []byte(fmt.Sprintf("%d",newval)))
	}

	return int(newval), nil
}

func (d *DynamoModule) Incr(key string) (int, error) {
	return d.Incrby(key, []byte("1"))
}
