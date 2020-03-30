package published

import (
	"gerrit-share.lan/go/actors"
	"log"
)

type publishedActorsActor struct {
	actors.Actor
	actors      actors.ActorStateChangeStream
	broadcaster actors.StateBroadcaster
}

func (p *publishedActorsActor) publish(actor actors.ActorService) {
	p.Monitor(actor, func(error) {
		p.actors.Remove(actor)
	})
	p.actors.Add(actor)
	p.broadcaster.NewDataAvailable()
}

func (p *publishedActorsActor) subscribe(command *subscribe) {
	p.InitStreamOutput(p.broadcaster.AddOutput(), command)
}

func (p *publishedActorsActor) Shutdown() error {
	log.Println("publisher plugin shut down")
	return nil
}

func (p *publishedActorsActor) MakeBehaviour() actors.Behaviour {
	log.Println("publisher plugin started")
	p.broadcaster = actors.NewBroadcaster(&p.actors)
	p.SetExitProcessor(func() { p.broadcaster.Close(nil) })
	var b actors.Behaviour
	b.AddCommand(new(publish), func(cmd interface{}) (actors.Response, error) {
		p.publish(cmd.(*publish).actor)
		return nil, nil
	}).AddCommand(new(subscribe), func(cmd interface{}) (actors.Response, error) {
		p.subscribe(cmd.(*subscribe))
		return nil, nil
	}).Result(new(actors.ActorsArray))
	return b
}

func init() {
	actors.AddPlugin(pluginName, func() actors.BehavioralActor { return new(publishedActorsActor) })
}

func init() {
}
