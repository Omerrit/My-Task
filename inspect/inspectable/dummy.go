package inspectable

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/inspect"
)

type DummyInspectable common.None

func (d *DummyInspectable) Inspect(inspector *inspect.GenericInspector) {

}
