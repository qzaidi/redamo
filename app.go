package main

import (
	redis "github.com/qzaidi/redamo/redis"
  store "github.com/qzaidi/redamo/store"
	logging "gopkg.in/tokopedia/logging.v1"
  "strings"
	"log"
	"os"
	"syscall"
)

// a simple mapper, given a key, it returns table name, key column, value col
// and the value of the key. our example has keys of the type
// tblname:keycol:valcol:keyname, e.g. shop_login:sid:followers:11212
func mapper(key string) (string,string,string,string) {
  vals := strings.Split(key,":")
  return vals[0],vals[1],vals[2],vals[3]
}

func main() {
	logging.LogInit()
	port := 6379

  store := store.NewDynamoModule(mapper)
	server, err := redis.NewRedamoServer(port,store)
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
