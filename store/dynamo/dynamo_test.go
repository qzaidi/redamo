package dynamo

import (
	"flag"
	"os"
	"testing"
)

var d *DynamoModule

func mapper(string) (tbl string, kcol string, vcol string, kval string) {
	return "shop_login", "sid", "testing", "11212"
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
