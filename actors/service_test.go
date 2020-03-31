package actors

import (
	"fmt"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/inspect/inspectables"
	"testing"
)

const intValue = 10

type dummyPing int

const dummyPingName = "dummy_ping"

func (d *dummyPing) Inspect(inspector *inspect.GenericInspector) {
	inspector.Int((*int)(d), dummyPingName, "")
}

func (*dummyPing) Fill() {}

func newDummyPing(value dummyPing) *dummyPing {
	return &value
}

type dummyPong int

type dummyData float64

func newDummyData(value dummyData) *dummyData {
	return &value
}

func (d *dummyData) Inspect(inspector *inspect.GenericInspector) {
	inspector.Float64((*float64)(d), 'f', 16, "dummyData", "")
}

type dummyService struct {
	Actor
	gotCommand bool
}

func (d dummyPong) Visit(visitor ResponseVisitor) {
	visitor.Reply(&d)
}

func (d *dummyPong) Inspect(inspector *inspect.GenericInspector) {
	inspector.Int((*int)(d), "dummyPong", "")
}

func newDummyPong(value dummyPong) *dummyPong {
	return &value
}

func (d *dummyService) MakeBehaviour() Behaviour {
	var b Behaviour
	b.AddCommand(new(dummyPing), func(command interface{}) (Response, error) {
		ping := command.(*dummyPing)
		d.gotCommand = *ping == intValue
		if d.gotCommand {
			fmt.Println("got command")
			return dummyPong(intValue * int(*ping)), nil
		} else {
			return nil, errors.New("Invalid value")
		}
	})
	return b
}

func newDummyService(t *testing.T) *dummyService {
	service := new(dummyService)
	printOnPanic(t, &service.Actor)
	return service
}

func printOnPanic(t *testing.T, actor *Actor) {
	actor.SetPanicProcessor(func(err errors.StackTraceError) {
		fmt.Println("panic:", errors.FullInfo(err))
		t.Error("panic: ", errors.FullInfo(err))
		actor.Quit(err)
	})
}

func ensureDead(t *testing.T, service ActorService, expectedErrors ...error) {
	service.System().Become(NewSimpleActor(func(actor *Actor) Behaviour {
		actor.Monitor(service, func(err error) {
			for _, e := range expectedErrors {
				if !errors.Is(err, e) {
					t.Error("actor closed with unexpected error:", errors.FullInfo(err))
				}
			}
			//actor should quit by itself
		})
		return Behaviour{}
	}))
}

func newTestingActor(t *testing.T, behaviourMaker BehaviourMaker) BehavioralActor {
	return NewSimpleActor(func(actor *Actor) Behaviour {
		printOnPanic(t, actor)
		if behaviourMaker != nil {
			return behaviourMaker(actor)
		}
		return Behaviour{}
	})
}

func init() {
	inspectables.Register(dummyPingName, func() inspect.Inspectable { return new(dummyPing) })
}

func TestProcessors(t *testing.T) {
	var s System
	s.EnableSurveillance()
	defer s.WaitFinished()
	service := s.Spawn(newDummyService(t))
	sender := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		actor.Depend(service)
		actor.SendRequest(service, newDummyPing(intValue), OnReply(func(reply interface{}) {
			fmt.Println("got reply")
			result := reply.(*dummyPong)
			if *result != intValue*intValue {
				t.Error("Bad response")
				return
			} else {
				t.Log(result)
			}
		}).OnError(actor.Quit))
		return Behaviour{}
	}))
	ensureDead(t, service, nil)
	ensureDead(t, sender, nil)
}

func TestProcessingError(t *testing.T) {
	var s System
	defer s.WaitFinished()
	service := s.Spawn(newDummyService(t))
	var err error
	sender := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		actor.Depend(service)
		actor.SendRequest(service, newDummyPing(intValue*2), OnReplyError(func(e error) {
			err = e
			actor.Quit(nil)
		}))
		return Behaviour{}
	}))
	if err == nil {
		t.Error("Should get error here")
		return
	}
	t.Log(err)
	ensureDead(t, service, nil)
	ensureDead(t, sender, nil)
}

func TestMonitorDeath(t *testing.T) {
	var s System
	defer s.WaitFinished()
	var count int
	monitorable := s.Spawn(NewSimpleActor(nil))
	monitor := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		actor.Monitor(monitorable, func(error) {
			count++
		})
		return Behaviour{}
	}))
	if count != 1 {
		t.Error("Wrong dead count")
	}
	ensureDead(t, monitorable, nil)
	ensureDead(t, monitor, nil)
}

/*
func TestRunAfter(t *testing.T) {
	var delta time.Duration
	service := MakeServiceAndStart(nil, StarterCommFunc(func(c *Communicator) {
		now := time.Now()
		c.SetCommandProcessor(CommandProcessorFunc(func(cmd interface{}) (Response, error) {
			if _, ok := cmd.(dummyPing); ok {
				delta = time.Now().Sub(now)
			} else {
				t.Error("bad command")
			}
			c.Quit(nil)
			return nil, nil
		}))
		c.RunCommandAfter(time.Second, dummyPing(0), nil)
	}))
	service.DoneChannel().Await()
	if delta < time.Second {
		t.Error("too soon")
	}
	fwtest.CloseService(t, service, nil)
}*/

func TestDelegate(t *testing.T) {
	var s System
	defer s.WaitFinished()
	service := s.Spawn(newDummyService(t))
	delegator := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		b.AddCommand(new(dummyPing), func(interface{}) (Response, error) {
			actor.Delegate(service)
			return nil, nil
		})
		actor.Link(service)
		return b
	}))
	requestor := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		actor.Link(delegator)
		actor.SendRequest(delegator, newDummyPing(intValue), OnReply(func(reply interface{}) {
			if *(reply.(*dummyPong)) == dummyPong(intValue*intValue) {
				actor.Quit(nil)
			} else {
				actor.Quit(errors.New("wrong value"))
			}
		}).QuitOnError(actor, "request failed"))
		actor.Link(delegator)
		return Behaviour{}
	}))
	ensureDead(t, requestor, nil)
	ensureDead(t, delegator, nil)
	ensureDead(t, service, nil)
}

func TestPause(t *testing.T) {
	var s System
	service := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		b.AddCommand(new(dummyPing), func(input interface{}) (Response, error) {
			return newDummyPong(dummyPong(*input.(*dummyPing)) + 1), nil
		})
		return b
	}))
	filter := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		actor.Link(service)
		b.PushCommandFilter(func(input interface{}) error {
			if ping, ok := input.(*dummyPing); ok {
				cmd := actor.PauseCommand()
				if cmd == nil {
					t.Error("failed to pause")
					actor.Quit(nil)
					return nil
				}
				actor.SendRequest(service, ping, OnReply(func(reply interface{}) {
					*ping = dummyPing(*(reply.(*dummyPong)))
					actor.ResumeCommand(cmd)
				}).OnError(func(err error) {
					actor.CancelCommand(cmd, err)
					actor.Quit(err)
				}))
			}
			return nil
		})
		b.AddCommand(new(dummyPing), func(input interface{}) (Response, error) {
			return newDummyPong(dummyPong(*input.(*dummyPing)) + 1), nil
		})
		return b
	}))
	requestor := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		actor.Link(filter)
		actor.SendRequest(filter, newDummyPing(intValue), OnReply(func(reply interface{}) {
			value := *(reply.(*dummyPong))
			if value != intValue+2 {
				t.Error("value is ", value, "but should be", intValue+2)
			} else {
				t.Log("ok")
			}
		}).QuitOnError(actor, "error sending initial ping request"))
		return Behaviour{}
	}))
	ensureDead(t, service, nil)
	ensureDead(t, filter, nil)
	ensureDead(t, requestor, nil)
}

type simpleRequestStream struct {
	RequestStreamBase
}

func init() {
	inspectables.Register("simpleRequestStream", func() inspect.Inspectable { return new(simpleRequestStream) })
}

func Test3Stream(t *testing.T) {
	var s System
	defer s.WaitFinished()
	var value1, value2 dummyData

	source := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		b.AddCommand(new(simpleRequestStream), func(command interface{}) (Response, error) {
			cmd := command.(*simpleRequestStream)
			fmt.Printf("%p source got init request\n", actor.Service())
			source := new(genericCallbackStreamOutput)
			sent := false
			source.generator = func(out *StreamOutputBase, data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
				if maxLen != 1 {
					actor.Quit(errors.New("source: Wrong request length"))
					return nil, nil
				}
				if _, ok := data.(*dummyData); !ok {
					actor.Quit(errors.New("source: wrong type requested"))
					return nil, nil
				}
				source.CloseWhenAcknowledged()
				fmt.Println("generating data")
				if !sent {
					sent = true
					return newDummyData(intValue), nil
				} else {
					return nil, nil
				}
			}
			actor.InitStreamOutput(source, cmd)
			actor.Quit(nil)
			return NoneReply{}, nil
		})
		return b
	}))

	delegator := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		b.AddCommand(new(simpleRequestStream), func(command interface{}) (Response, error) {
			actor.Delegate(source)
			return nil, nil
		})
		actor.DependOn(source)
		return b
	}))

	processor := s.Spawn(newTestingActor(t, func(actor *Actor) Behaviour {
		var b Behaviour
		b.AddCommand(new(simpleRequestStream), func(command interface{}) (Response, error) {
			cmd := command.(*simpleRequestStream)
			fmt.Printf("%p :processor got init request\n", actor.Service())
			output := new(genericCallbackStreamOutput)
			var value *dummyData
			output.generator = func(base *StreamOutputBase, data inspect.Inspectable, maxLen int) (inspect.Inspectable, error) {
				if maxLen != 1 {
					actor.Quit(errors.New("processor: Wrong request length"))
					return nil, nil
				}
				if _, ok := data.(*dummyData); !ok {
					actor.Quit(errors.New("processor: wrong type requested"))
					return nil, nil
				}
				if value != nil {
					v := *value
					output.CloseWhenAcknowledged()
					value = nil
					return &v, nil
				}
				return nil, nil
			}
			input := NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
				if data == nil {
					return nil
				}
				if dummy, ok := data.(*dummyData); ok {
					value = dummy
					value1 = *dummy
					output.FlushLater()
					return nil
				}
				actor.Quit(errors.New("processor: got wrong type"))
				return nil
			}, func(base *StreamInputBase) {
				base.RequestData(newDummyData(0), 1)
			})
			input.OnClose(output.CloseStream)
			actor.InitStreamOutput(output, cmd)
			actor.RequestStream(input, delegator, new(simpleRequestStream), actor.Quit)
			actor.Quit(nil)
			return NoneReply{}, nil
		})
		return b
	}))

	sink := s.Become(newTestingActor(t, func(actor *Actor) Behaviour {
		sink := NewSimpleCallbackStreamInput(func(data inspect.Inspectable) error {
			fmt.Printf("%p sink got data\n", actor.Service())
			if data == nil {
				return nil
			}
			if value, ok := data.(*dummyData); ok {
				value2 = *value
				return nil
			}
			actor.Quit(errors.New("sink:got wrong type"))
			return nil
		}, func(base *StreamInputBase) {
			base.RequestData(newDummyData(0), 1)
		})
		sink.OnClose(func(err error) {
			fmt.Println("Sink closed: ", err)
		})
		actor.RequestStream(sink, processor, new(simpleRequestStream), actor.Quit)
		return Behaviour{}
	}))

	ensureDead(t, sink, nil)
	ensureDead(t, source, nil)
	ensureDead(t, processor, nil)
	ensureDead(t, delegator, nil)

	if value1 != dummyData(intValue) {
		t.Error("wrong value 1")
	}
	if value2 != dummyData(intValue) {
		t.Error("wrong value 2")
	}
}
