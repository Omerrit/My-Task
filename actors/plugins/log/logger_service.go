package log

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/utils/flags"
	"log"
)

const (
	pluginName        = "logger"
	defaultStackDepth = 1024
)

type loggerService struct {
	actors.Actor
	broadcaster actors.StateBroadcaster
	stackDepth  uint
	messages    MessagesStream
	writers     map[int]common.None
}

func (l *loggerService) MakeBehaviour() actors.Behaviour {
	log.Println("logger started")
	l.initWriters()
	var b actors.Behaviour
	l.stackDepth = defaultStackDepth
	l.broadcaster = actors.NewBroadcaster(&l.messages)
	b.AddCommand(new(subscribe), func(cmd interface{}) (actors.Response, error) {
		l.subscribe(cmd.(*subscribe))
		return nil, nil
	})
	b.AddMessage(new(notGonnaSubscribe), func(cmd interface{}) {
		notGonnaCmd := cmd.(*notGonnaSubscribe)
		if l.shouldStartDeleting(notGonnaCmd.id) {
			l.messages.deleteHistory = true
		}
	})
	b.AddMessage(new(Message), func(cmd interface{}) {
		l.writeMessage(cmd.(*Message))
	})

	l.SetPanicProcessor(func(err errors.StackTraceError) {
		fmt.Println(err.StackTrace())
		l.Quit(err)
	})
	return b
}

func (l *loggerService) Shutdown() error {
	log.Println("logger shut down")
	return nil
}

func (l *loggerService) writeMessage(command *Message) {
	l.messages.Add(*command)
	l.broadcaster.NewDataAvailable()
	return
}

func (l *loggerService) subscribe(command *subscribe) {
	l.InitStreamOutput(l.broadcaster.AddOutput(), command)
	if l.shouldStartDeleting(command.id) {
		l.messages.deleteHistory = true
	}
}

func (l *loggerService) initWriters() {
	l.writers = make(map[int]common.None, writerCounter)
	for i := 0; i < int(writerCounter); i++ {
		l.writers[i] = common.None{}
	}
}

func (l *loggerService) shouldStartDeleting(id int) bool {
	if _, ok := l.writers[id]; !ok {
		log.Println("unexpected writers id detected")
		return false
	}
	delete(l.writers, id)
	if len(l.writers) > 0 {
		return false
	}
	return true
}

func init() {
	actors.AddPlugin(pluginName, func() actors.BehavioralActor { return new(loggerService) })

	starter.SetFlagInitializer(pluginName, func() {
		flags.IntFlag(&verbosityFlag, "verb", "logger verbosity")
	})
}
