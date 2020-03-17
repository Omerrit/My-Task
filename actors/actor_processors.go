package actors

import (
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/mreflect"
)

type actorProcessors struct {
	commandProcessors        compiledCommandBehaviour
	messageProcessors        compiledMessageBehaviour
	defaultCommandProcessors compiledCommandBehaviour
	defaultMessageProcessors compiledMessageBehaviour
	finishedServiceProcessor FinishedServiceProcessor
	panicProcessor           PanicProcessor
	exitProcessor            common.SimpleCallback
	messageProcessorsCleared bool

	currentFilterIndex int
	currentCommand     commandMessage
}

func (a *actorProcessors) clearMessageProcessors() {
	a.commandProcessors = a.defaultCommandProcessors
	a.messageProcessors = a.defaultMessageProcessors
	a.panicProcessor = nil
	a.finishedServiceProcessor = nil
	a.messageProcessorsCleared = true
}

func (a *actorProcessors) haveActiveProcessors() bool {
	return !a.messageProcessorsCleared
}

//default behaviour can be overriden
func (a *actorProcessors) setDefaultBehaviour(system *System, behaviour Behaviour) {
	a.defaultCommandProcessors = compileCommandBehaviour(system, behaviour.CommandBehaviour, "", nil)
	a.defaultMessageProcessors = compileMessageBehaviour(system, behaviour.MessageBehaviour)
}

func (a *actorProcessors) setBehaviour(system *System, behaviour Behaviour) {
	commandProcessors := compileCommandBehaviour(system, behaviour.CommandBehaviour, behaviour.Name, behaviour.CommandFilters)
	messageProcessors := compileMessageBehaviour(system, behaviour.MessageBehaviour)
	//don't set processors if compile* panics
	a.messageProcessorsCleared = commandProcessors.IsEmpty() && messageProcessors.IsEmpty()
	a.commandProcessors = *commandProcessors.addNewBindings(&a.defaultCommandProcessors)
	a.messageProcessors = *messageProcessors.addNewBindings(&a.defaultMessageProcessors)
}

func (a *actorProcessors) SetFinishedServiceProcessor(processor FinishedServiceProcessor) {
	a.finishedServiceProcessor = processor
}

func (a *actorProcessors) SetPanicProcessor(processor PanicProcessor) {
	a.panicProcessor = processor
}

func (a *actorProcessors) SetExitProcessor(processor common.SimpleCallback) {
	a.exitProcessor = processor
}

//filter #0 is the command processor
func (a *actorProcessors) runCommandProcessor(cmd interface{}, filterIndex int) (Response, error) {
	processor := a.commandProcessors.processors[mreflect.GetTypeId(cmd)]
	if processor != nil {
		var err error
		if filterIndex > 0 {
			if filterIndex > len(a.commandProcessors.commandFilters) {
				a.currentFilterIndex = len(a.commandProcessors.commandFilters) - 1
			} else {
				a.currentFilterIndex = filterIndex - 1
			}
			for ; a.currentFilterIndex >= 0; a.currentFilterIndex-- {
				err = a.commandProcessors.commandFilters[a.currentFilterIndex](a.currentCommand.data)
				if err != nil {
					return nil, err
				}
				if !a.currentCommand.isValid() {
					return nil, nil
				}
			}
		}
		return processor(cmd)
	}
	return nil, ErrUnknownCommand
}

func (a *actorProcessors) runMessageProcessor(msg interface{}) {
	processor := a.messageProcessors.processors[mreflect.GetTypeId(msg)]
	if processor != nil {
		processor(msg)
	}
}

func (a *actorProcessors) runFinishedServiceProcessor(service ActorService, err error) {
	if a.finishedServiceProcessor != nil {
		a.finishedServiceProcessor(service, err)
	}
}

func (a *actorProcessors) runPanicProcessor(err errors.StackTraceError) {
	if a.panicProcessor != nil {
		a.panicProcessor(err)
	}
}

func (a *actorProcessors) runExitProcessor() {
	if a.exitProcessor != nil {
		a.exitProcessor()
		a.exitProcessor = nil
	}
}
