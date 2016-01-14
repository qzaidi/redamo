package main

import (
	redis "github.com/qzaidi/redamo/redis"
	store "github.com/qzaidi/redamo/store/dynamo"
	logging "gopkg.in/tokopedia/logging.v1"
	"log"
	"os"
	"strings"
	"syscall"
)

// a simple mapper, given a key, it returns table name, key column, value col
// and the value of the key. our example has keys of the type
// tblname:keycol:valcol:keyname, e.g. shop_login:sid:followers:11212
func mapper(key string) *store.Mapper {

  // defaults
  kmap := &store.Mapper {
    Table: "redamo",
    Kcol: "key",
    Vcol: "val",
    Keyval: "12345",
  }

	vals := strings.Split(key, ":")

  if len(vals) == 4 {
    kmap.Table = vals[0]
    kmap.Kcol = vals[1]
    kmap.Vcol = vals[2]
    kmap.Keyval = vals[3]
  }
	return kmap
}

func main() {
	logging.LogInit()
	port := 6379

	dyn := store.NewDynamoModule(nil)
	server, err := redis.NewRedamoServer(port, dyn)
	if err != nil {
		panic(err)
	}

	log.Println("Redamo listening on port", port)

	if err := server.ListenAndServe(); err != nil {
		if syserr, ok := err.(syscall.Errno); ok {
			if syserr == syscall.EADDRINUSE {
				log.Println(syserr)
				os.Exit(1)
			}
		}
		panic(err)
	}
}
