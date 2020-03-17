package actors

import (
	"gerrit-share.lan/go/inspect"
)

type commandBinding struct {
	command      inspect.Inspectable
	resultSample inspect.Inspectable
	resultType   inspect.TypeId
	processor    CommandProcessor
}

type CommandBehaviour []commandBinding

func (c *CommandBehaviour) AddCommand(command inspect.Inspectable, processor CommandProcessor) *CommandBehaviour {
	*c = append(*c, commandBinding{command, nil, inspect.TypeInvalid, processor})
	return c
}

func (c *CommandBehaviour) setLast(setter func(*commandBinding)) *CommandBehaviour {
	if len(*c) == 0 {
		return c
	}
	setter(&(*c)[len(*c)-1])
	return c
}

//sets result (type and sample) of the last command to its argument
func (c *CommandBehaviour) Result(sample inspect.Inspectable) *CommandBehaviour {
	return c.setLast(func(b *commandBinding) {
		b.resultSample = sample
		b.resultType = inspect.TypeValue
	})
}

func (c *CommandBehaviour) setResultType(resultType inspect.TypeId) *CommandBehaviour {
	return c.setLast(func(b *commandBinding) {
		b.resultType = resultType
	})
}

//sets result type of the last command to int
func (c *CommandBehaviour) ResultInt() *CommandBehaviour {
	return c.setResultType(inspect.TypeInt)
}

//sets result type of the last command to int32
func (c *CommandBehaviour) ResultInt32() *CommandBehaviour {
	return c.setResultType(inspect.TypeInt32)
}

//sets result type of the last command to int64
func (c *CommandBehaviour) ResultInt64() *CommandBehaviour {
	return c.setResultType(inspect.TypeInt64)
}

//sets result type of the last command to float32
func (c *CommandBehaviour) ResultFloat32() *CommandBehaviour {
	return c.setResultType(inspect.TypeFloat32)
}

//sets result type of the last command to float64
func (c *CommandBehaviour) ResultFloat64() *CommandBehaviour {
	return c.setResultType(inspect.TypeFloat64)
}

//sets result type of the last command to string
func (c *CommandBehaviour) ResultString() *CommandBehaviour {
	return c.setResultType(inspect.TypeString)
}

//sets result type of the last command to []byte that should be interpreted as string
func (c *CommandBehaviour) ResultByteString() *CommandBehaviour {
	return c.setResultType(inspect.TypeByteString)
}

//sets result type of the last command to []byte
func (c *CommandBehaviour) ResultBytes() *CommandBehaviour {
	return c.setResultType(inspect.TypeBytes)
}

//sets result type of the last command to bool
func (c *CommandBehaviour) ResultBool() *CommandBehaviour {
	return c.setResultType(inspect.TypeBool)
}

//sets result type of the last command to *big.Int
func (c *CommandBehaviour) ResultBigInt() *CommandBehaviour {
	return c.setResultType(inspect.TypeBigInt)
}

//sets result type of the last command to *big.Rat
func (c *CommandBehaviour) ResultRat() *CommandBehaviour {
	return c.setResultType(inspect.TypeRat)
}

//sets result type of the last command to *big.Float
func (c *CommandBehaviour) ResultBigFloat() *CommandBehaviour {
	return c.setResultType(inspect.TypeBigFloat)
}

type messageBinding struct {
	message   inspect.Inspectable
	processor MessageProcessor
}
type MessageBehaviour []messageBinding

func (m *MessageBehaviour) AddMessage(message inspect.Inspectable, processor MessageProcessor) *MessageBehaviour {
	*m = append(*m, messageBinding{message, processor})
	return m
}

type CommandFilter func(interface{}) error

type CommandFilters []CommandFilter

func (c *CommandFilters) PushCommandFilter(filter CommandFilter) {
	*c = append(*c, filter)
}

type Behaviour struct {
	CommandBehaviour
	MessageBehaviour
	CommandFilters
	Name string
}

//TODO: make actor return behaviour and set it when actor starts
