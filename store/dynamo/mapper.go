package dynamo

import (
	"gopkg.in/tokopedia/logging.v1"
	"strings"
)

func (d *DynamoModule) defaultMapper(key string) *Mapper {

	debug := logging.Debug.Println

	m := d.config.Keymap
	var kd *Mapper

	for k := range m {
		if strings.HasPrefix(key, k) {
			kd = &Mapper{}

			kd.Keyval = strings.TrimPrefix(key, k)
			kd.Table = m[k].Table
			kd.Kcol = m[k].Kcol
			kd.Vcol = m[k].Vcol
			kd.Ktype = m[k].Ktype
			kd.Vtype = m[k].Vtype
			debug("match", k, m[k])
			break
		}
	}

	return kd
}
