package auth

import (
	"gerrit-share.lan/go/inspect"
)

type Id [16]byte //uuid is enough

var zeroId Id

const IdName = packageName + ".id"

func (c *Id) Inspect(inspector *inspect.GenericInspector) {
	const description = "128 bit identifier"
	if inspector.IsReading() {
		var in []byte
		inspector.Bytes(&in, IdName, description)
		if len(in) != 16 {
			inspector.SetError(ErrWrongIdLen)
			return
		}
		copy((*c)[:], in)
	} else {
		out := (*c)[:]
		inspector.Bytes(&out, IdName, description)
	}
}

//sql.Scanner compatibility
func (c *Id) Scan(data interface{}) error {
	bytes, ok := data.([]byte)
	if !ok {
		return ErrWrongIdType
	}
	if len(bytes) != 16 {
		return ErrWrongIdLen
	}
	copy((*c)[:], bytes)
	return nil
}

type IdCallback func(id Id)

func (c IdCallback) Call(id Id) {
	if c != nil {
		c(id)
	}
}
