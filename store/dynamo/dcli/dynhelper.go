//helper for creating dynamo tables 
package main

import (
  "fmt"
  "flag"
  "os"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
)

func createTable(endpoint string,name string, pkey string, rc int64, wc int64) error {

  config := &aws.Config{
               Region: aws.String("ap-southeast-1"),
               DisableSSL: aws.Bool(true),
            }

  if endpoint != "" {
    config.Endpoint = aws.String(endpoint)
  }

  sess := session.New(config)

  client := dynamodb.New(sess)

  tablespecs := &dynamodb.CreateTableInput{
    TableName: aws.String(name),
    KeySchema: []*dynamodb.KeySchemaElement{
      {
        AttributeName: aws.String(pkey),
        KeyType: aws.String(dynamodb.KeyTypeHash),
      },
    },
    ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
      ReadCapacityUnits: aws.Int64(rc),
      WriteCapacityUnits: aws.Int64(wc),
    },
    AttributeDefinitions: []*dynamodb.AttributeDefinition{
        {
            AttributeName: aws.String(pkey),
            AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
        },
    },
  }
  _, err := client.CreateTable(tablespecs)
  return err
}

func main() {
  endpoint := flag.String("endpoint","","dynamodb endpoint, localhost:4567 for dynalite")
  tbl := flag.String("table", "REQUIRED", "create named table")
  rc :=  flag.Int64("rc",1,"read capacity")
  wc :=  flag.Int64("wc",1,"write capacity")
  key := flag.String("key","REQUIRED", "primary key for new table")

  flag.Parse()
  if *tbl == "REQUIRED" || *key == "REQUIRED" {
    fmt.Println("missing parameters, table and key is mandatory")
    os.Exit(1)
  }
  err := createTable(*endpoint,*tbl,*key,*rc,*wc)
  if err != nil {
    fmt.Println(err)
  }
}
