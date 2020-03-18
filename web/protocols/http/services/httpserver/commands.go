package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
)

type commandBatch struct {
	commands  []inspect.Inspectable
	generator inspectables.Creator
}

const commandBatchName = packageName + ".commandbatch"

func (c *commandBatch) Inspect(i *inspect.GenericInspector) {
	arrayInspector := i.Array(commandBatchName, "", "")
	{
		if !arrayInspector.IsReading() {
			return
		}
		c.commands = make([]inspect.Inspectable, arrayInspector.GetLength())
		for i := 0; i < len(c.commands); i++ {
			c.commands[i] = c.generator()
		}
		for index := range c.commands {
			c.commands[index].Inspect(arrayInspector.Value())
		}
		arrayInspector.End()
	}
}
