package dynamo

import (
	"strings"
  "log"
)

func (d *DynamoModule) defaultMapper(key string) *Mapper {

	m := d.config.Keymap
  var kd *Mapper

  log.Println(m)

	for k := range m {
    log.Println("checking",k)
		if strings.HasPrefix(key, k) {
      kd = &Mapper{}

			kd.Keyval = strings.TrimPrefix(key, k)
			kd.Table = m[k].Table
			kd.Kcol = m[k].Kcol
			kd.Vcol = m[k].Vcol
      kd.Ktype = m[k].Ktype
      kd.Vtype = m[k].Vtype
			break
		}
	}

	return kd
}
