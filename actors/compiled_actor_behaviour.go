package actors

import (
	"fmt"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectwrappers"
	"gerrit-share.lan/go/mreflect"
	"reflect"
)

type ActorCommand struct {
	Command    PortableType
	Result     PortableType
	ResultType inspect.TypeId
}

const ActorCommandName = packageName + ".actor_command"

type ActorCommands struct {
	Commands []ActorCommand
	Name     string
}

const ActorCommandsName = packageName + ".actor_commands"

//TODO: make ActorCommands inspectable
//how should we load and save unknown inspectables?
//probably save as type name+data
//and load as type name and then get concrete type
//it should be provided by embedded loader
//that can also spawn actors

type compiledCommandBehaviour struct {
	commands       ActorCommands
	processors     map[mreflect.TypeId]CommandProcessor
	commandFilters CommandFilters
}

func makeActorCommand(system *System, commandName string, binding commandBinding) ActorCommand {
	if binding.resultType == inspect.TypeValue {
		resultName, ok := system.types.GetName(binding.resultSample)
		if !ok {
			panic(fmt.Sprintf("result %s is not registered", reflect.TypeOf(binding.resultSample)))
		}
		return ActorCommand{PortableType{commandName, binding.command},
			PortableType{resultName, binding.resultSample},
			binding.resultType}
	} else if binding.resultType == inspect.TypeInvalid {
		return ActorCommand{PortableType{commandName, binding.command}, PortableType{}, binding.resultType}
	} else {
		value, err := inspectwrappers.NewDefaultBasicValue(binding.resultType)
		if err != nil {
			panic(fmt.Errorf("%w: %d", err, binding.resultType))
		}
		return ActorCommand{PortableType{commandName, binding.command}, PortableType{"", value}, binding.resultType}
	}
}

func compileCommandBehaviour(system *System, behaviour CommandBehaviour, name string, filters CommandFilters) compiledCommandBehaviour {
	var compiled compiledCommandBehaviour
	compiled.commands.Name = name
	if len(behaviour) == 0 {
		return compiled
	}
	compiled.commandFilters = filters
	compiled.processors = make(map[mreflect.TypeId]CommandProcessor, len(behaviour))
	compiled.commands.Commands = make([]ActorCommand, 0, len(behaviour))
	indices := make(map[mreflect.TypeId]int, len(behaviour))
	for _, binding := range behaviour {
		typeid := mreflect.GetTypeId(binding.command)
		name, ok := system.types.GetNameById(typeid)
		if !ok {
			panic(fmt.Sprintf("command %s is not registered", reflect.TypeOf(binding.command)))
		}
		command := makeActorCommand(system, name, binding)
		if index, ok := indices[typeid]; ok {
			compiled.commands.Commands[index] = command
		} else {
			indices[typeid] = len(compiled.commands.Commands)
			compiled.commands.Commands = append(compiled.commands.Commands, command)
		}
		compiled.processors[typeid] = binding.processor
	}
	return compiled
}

func (c *compiledCommandBehaviour) Clear() {
	c.commands.Commands = nil
	c.processors = nil
}

func (c *compiledCommandBehaviour) IsEmpty() bool {
	return c.commands.Commands == nil && c.processors == nil
}

func (c *compiledCommandBehaviour) addDefaultBindings(system *System, behaviour CommandBehaviour) {
	for _, binding := range behaviour {
		typeid := mreflect.GetTypeId(binding.command)
		name, ok := system.types.GetNameById(typeid)
		if !ok {
			panic(fmt.Sprintf("command %s is not registered", reflect.TypeOf(binding.command)))
		}
		if _, ok := c.processors[typeid]; ok {
			c.commands.Commands = append(c.commands.Commands, makeActorCommand(system, name, binding))
			c.processors[typeid] = binding.processor
		}
	}
}

func (c *compiledCommandBehaviour) addNewBindings(other *compiledCommandBehaviour) *compiledCommandBehaviour {
	if c.commands.Commands == nil || c.processors == nil {
		name := c.commands.Name
		*c = *other
		c.commands.Name = name
		return c
	}
	for _, command := range other.commands.Commands {
		typeid := mreflect.GetTypeId(command.Command.sample)
		if _, ok := c.processors[typeid]; !ok {
			c.commands.Commands = append(c.commands.Commands, command)
			c.processors[typeid] = other.processors[typeid]
		}
	}
	return c
}

type compiledMessageBehaviour struct {
	messages   PortableValues
	processors map[mreflect.TypeId]MessageProcessor
}

func compileMessageBehaviour(system *System, behaviour MessageBehaviour) compiledMessageBehaviour {
	var compiled compiledMessageBehaviour
	if len(behaviour) == 0 {
		return compiled
	}
	compiled.processors = make(map[mreflect.TypeId]MessageProcessor, len(behaviour))
	compiled.messages = make([]PortableType, 0, len(behaviour))
	indices := make(map[mreflect.TypeId]int, len(behaviour))
	for _, binding := range behaviour {
		typeid := mreflect.GetTypeId(binding.message)
		name, ok := system.types.GetNameById(typeid)
		if !ok {
			panic(fmt.Sprintf("message %s is not registered", reflect.TypeOf(binding.message)))
		}
		if index, ok := indices[typeid]; ok {
			compiled.messages[index] = PortableType{name, binding.message}
		} else {
			indices[typeid] = len(compiled.messages)
			compiled.messages = append(compiled.messages, PortableType{name, binding.message})
		}
		compiled.processors[typeid] = binding.processor
	}
	return compiled
}

func (c *compiledMessageBehaviour) Clear() {
	c.messages = nil
	c.processors = nil
}

func (c *compiledMessageBehaviour) IsEmpty() bool {
	return c.messages == nil && c.processors == nil
}

func (c *compiledMessageBehaviour) addNewBindings(other *compiledMessageBehaviour) *compiledMessageBehaviour {
	if c.processors == nil || c.messages == nil {
		*c = *other
		return c
	}
	for _, message := range other.messages {
		typeid := mreflect.GetTypeId(message.sample)
		if _, ok := c.processors[typeid]; !ok {
			c.messages = append(c.messages, message)
			c.processors[typeid] = other.processors[typeid]
		}
	}
	return c
}

func (a *ActorCommand) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(ActorCommandName, "actor command description")
	{
		a.Command.Inspect(o.Value("command", true, "command description"))
		a.ResultType.Inspect(o.Value("result_type", true, "result type"))
		if a.ResultType == inspect.TypeInvalid {
			if o.IsReading() {
				a.Result.sample = nil
			}
		} else if a.ResultType == inspect.TypeValue {
			a.Result.Inspect(o.Value("result", true, "result sample"))
		} else {
			if o.IsReading() {
				a.Result.sample, _ = inspectwrappers.NewDefaultBasicValue(a.ResultType)
			}
		}
		o.End()
	}
}

func (a *ActorCommands) Inspect(inspector *inspect.GenericInspector) {
	o := inspector.Object(ActorCommandsName, "actor commands")
	{
		o.String(&a.Name, "name", true, "service display name if available")
		i := o.Value("commands", true, "list of actor commands").Array("", ActorCommandName, "")
		{
			if i.IsReading() {
				length := i.GetLength()
				if cap(a.Commands) < length {
					a.Commands = make([]ActorCommand, length)
				} else {
					a.Commands = a.Commands[:length]
				}
			} else {
				i.SetLength(len(a.Commands))
			}
			for index := range a.Commands {
				a.Commands[index].Inspect(i.Value())
			}
			i.End()
		}
		o.End()
	}
}
