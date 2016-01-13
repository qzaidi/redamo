package dynamo

import (
	"flag"
	"os"
	"testing"
  "strings"
)

var d *DynamoModule

func mapper(k string) (tbl string, kcol string, vcol string, kval string) {
  f := strings.Split(k,":")
	return "shop_login", "sid", f[1], f[2] 
}

func TestMain(m *testing.M) {
	flag.Parse()
	d = NewDynamoModule(mapper)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	v, e := d.Get("s:testing:11212")
	t.Log(string(v), e)
}

func TestSet(t *testing.T) {
	e := d.Set("s:testing:11212", []byte("20"))
	t.Log(e)
}

func TestIncrby(t *testing.T) {
	val,e := d.Incrby("s:numeric:11212", 4)
	t.Log(val,e)
}
