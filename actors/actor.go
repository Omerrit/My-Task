package actors

import (
	"gerrit-share.lan/go/actors/internal/queue"
	"gerrit-share.lan/go/common"
	"gerrit-share.lan/go/debug"
	"gerrit-share.lan/go/errors"
	"gerrit-share.lan/go/inspect"
	"gerrit-share.lan/go/interfaces"
)

type ActorState int

const (
	ActorRunning ActorState = iota
	ActorQuitting
	ActorClosed
	ActorDead
)

type Actor struct {
	actorProcessors
	system           *System
	service          ActorService
	queue            *queue.Queue
	messages         []interface{}
	reissuedCommands reissuedQueue

	nextCommandId    commandId
	inflightRequests inflightRequests
	activePromises   promiseIdCallbacks

	links            links
	monitoringActors ActorErrorCallbacks
	awaitingActors   ActorErrorCallbacks

	currentStreamId streamId
	streamInputs    streamInputs
	streamOutputs   streamOutputs
	readyOutputs    StreamOutputSet

	state                  ActorState
	stateChangeSubscribers ActorSet
	quitError              error
}

//side note: command buffer for slow actors should support cancellation too and it's painful to find command for specific sender and id
//for a real queue (only by scanning), only ring buffer could make it quicker as every item have fixed absolute index

//actor system should also be passed here
//actor system shouldn't have additional modules, only basic functionality

func (a *Actor) init(system *System) {
	a.system = system
	a.setState(ActorRunning)
}

func (a *Actor) Service() ActorService {
	if a.service == nil {
		a.createService()
	}
	return a.service
}

func (a *Actor) System() *System {
	return a.system
}

func (a *Actor) GetBase() *Actor {
	return a
}

func (a *Actor) createService() {
	if a.queue == nil {
		a.queue = queue.NewQueue()
	}
	if a.service == nil {
		a.service = &service{}
		a.service.init(a, a.queue)
	}
}

func (a *Actor) setBehaviour(system *System, behaviour Behaviour) {
	defer func() {
		err := errors.RecoverToError(recover())
		if err != nil {
			a.onPanic(err)
		}
	}()
	var defaults Behaviour
	defaults.AddCommand(GetStatus{}, func(interface{}) (Response, error) {
		result := a.GetStatus()
		return &result, nil
	}).Result(new(Status))
	a.actorProcessors.setDefaultBehaviour(system, defaults)
	a.actorProcessors.setBehaviour(system, behaviour)
}

func (a *Actor) setState(state ActorState) {
	if a.state == state {
		return
	}
	a.state = state
	for actor := range a.stateChangeSubscribers {
		actor.enqueue(notifyStateChange{a.Service(), state})
	}
}

func (a *Actor) SendRequest(destination ActorService, request inspect.Inspectable, replyProcessor interfaces.ReplyProcessor) interfaces.Canceller {
	if a.state == ActorClosed {
		return nil
	}
	a.inflightRequests.Add(a.nextCommandId, replyProcessor)
	enqueue(destination, commandMessage{promiseId{a.Service(), a.nextCommandId}, request})
	a.nextCommandId++
	return &requestCanceller{a.nextCommandId - 1, a}
}

func (a *Actor) cancelRequest(id commandId, info requestInfo) {
	if info.destination != nil {
		enqueue(info.destination, cancelCommand{a.Service(), id})
	}
	info.processor.Error(ErrCancelled)
}

func (a *Actor) cancelRequestById(id commandId) {
	info, ok := a.inflightRequests[id]
	if ok {
		a.cancelRequest(id, info)
	}
}

func (a *Actor) SendMessage(destination ActorService, message interface{}) {
	enqueue(destination, message)
}

func (a *Actor) SendMessages(destination ActorService, messages []inspect.Inspectable) {
	enqueue(destination, messages)
}

//Close current service when destination closes
//and close other service when current service closes.
func (a *Actor) Link(destination ActorService) {
	a.links.Add(destination, linkLink)
	enqueue(destination, establishLink{a.Service(), linkLink})
}

func (a *Actor) monitor(destination ActorService) {
	enqueue(destination, establishLink{a.Service(), linkMonitor})
}

//callback would be called when destination would quit
//callback will never be called if this actor quits first
func (a *Actor) Monitor(destination ActorService, onExit common.ErrorCallback) {
	a.monitoringActors.Add(destination, onExit)
	a.monitor(destination)
}

//Close current service when destination closes
func (a *Actor) DependOn(destination ActorService) {
	enqueue(destination, establishLink{a.Service(), linkDepend})
}

//Close destination when current service closes.
func (a *Actor) Depend(destination ActorService) {
	a.links.Add(destination, linkKill)
}

//wait till destination quits then call callback
func (a *Actor) Await(destination ActorService, onExit common.ErrorCallback) {
	a.awaitingActors.Add(destination, onExit)
	a.monitor(destination)
}

func (a *Actor) MonitorStateChanges(destination ActorService) {
	enqueue(destination, subscribeStateChange{a.Service()})
}

func (a *Actor) close() {
	a.setState(ActorClosed)
	a.clearMessageProcessors()
	a.monitoringActors.Clear()
	for id := range a.activePromises {
		a.cancelCommandProcessingFromPromise(id, ErrActorDead)
	}
	for id, info := range a.inflightRequests {
		a.cancelRequest(id, info)
	}
	a.inflightRequests.Clear()
	for actor, linkType := range a.links {
		switch linkType {
		case linkLink, linkDepend, linkKill:
			actor.SendQuit(nil)
			a.monitor(actor)
		case linkMonitor:
			enqueue(actor, notifyClose{a.Service(), a.quitError})
			a.links.Remove(actor)
		}
	}
	for _, input := range a.streamInputs {
		a.closeInput(input, a.quitError)
	}
	for _, info := range a.streamOutputs {
		a.closeOutput(info.output, a.quitError)
	}
	a.readyOutputs.Clear()
	if !a.links.IsEmpty() {
		a.waitDependentActors()
	}
	a.setState(ActorDead)
	a.stateChangeSubscribers.Clear()
}

func (a *Actor) closeStreamOuts() {
	for _, out := range a.streamOutputs {
		base := out.output.getBase()
		if base.shouldCloseWhenActorCloses {
			//			fmt.Println("closing stream")
			base.CloseStream(a.quitError)
		} else {
			//			fmt.Println("not closing stream")
		}
	}
}

func (a *Actor) closeStreamIns() {
	for _, in := range a.streamInputs {
		base := in.getBase()
		if base.shouldCloseWhenActorCloses {

			base.Close()
		}
	}
}

func (a *Actor) Quit(err error) {
	a.quitError = err
	if a.state != ActorRunning {
		//		fmt.Println("Quit while not running, active processors:", a.haveActiveProcessors())
		return
	}
	a.setState(ActorQuitting)
	a.clearMessageProcessors()
	a.closeStreamOuts()
	a.closeStreamIns()
}

func (a *Actor) reply(data interface{}) {
	if a.currentCommand.isValid() {
		a.currentCommand.reply(data)
		a.currentCommand.invalidate()
	}
}

func (a *Actor) replyWithError(err error) {
	if a.currentCommand.isValid() {
		a.currentCommand.replyWithError(err)
		a.currentCommand.invalidate()
	}
}

func (a *Actor) promiseReply(id promiseId, data interface{}) {
	if a.activePromises.Contains(id) {
		id.reply(data)
		a.activePromises.Delete(id)
	}
}

func (a *Actor) promiseReplyWithError(id promiseId, err error) {
	if a.activePromises.Contains(id) {
		id.reply(err)
		a.activePromises.Delete(id)
	}
}

func (a *Actor) Delegate(destination ActorService) {
	if a.currentCommand.isValid() {
		if a.activePromises.Contains(a.currentCommand.promiseId) {
			//handle delegation when reprocessing command
			//source actor still thinks it can send cancel request here
			//but we can't properly process it since we're not processing a command
			//we or the destination can send preReply but it still takes time for a source
			//to understand that destination had changed

			//do proxying for now, it should also work with stream commands messages as there are separate stream messages
			id := a.currentCommand.promiseId
			canceler := a.SendRequest(destination, a.currentCommand.data, OnReply(func(data interface{}) {
				id.reply(data)
				a.activePromises.Delete(id)
			}).OnError(func(err error) {
				id.replyWithError(err)
				a.activePromises.Delete(id)
			}))
			a.activePromises.Add(id, canceler.Cancel)
			a.currentCommand.invalidate()
			return
		}
		enqueue(destination, a.currentCommand)
		a.currentCommand.invalidate()
	}
}

func (a *Actor) PauseCommand() *Command {
	if !a.currentCommand.isValid() {
		return nil
	}
	cmd := &Command{resumedCommand{a.currentCommand, a.currentFilterIndex}}
	a.activePromises.Add(a.currentCommand.promiseId, cmd.invalidate)
	a.currentCommand.invalidate()
	return cmd
}

func (a *Actor) PauseCommandAndEnqueue(queue *CommandQueue) {
	if !a.currentCommand.isValid() {
		return
	}
	id := a.currentCommand.promiseId
	queue.push(resumedCommand{a.currentCommand, a.currentFilterIndex})
	a.activePromises.Add(id, func() {
		queue.addCanceled(id)
	})
	a.currentCommand.invalidate()
}

func (a *Actor) ResumeCommand(queue commandPopper) {
	a.reissuedCommands.pushOne(queue)
}

func (a *Actor) ResumeCommands(queue commandPopper) {
	a.reissuedCommands.pushAll(queue)
}

func (a *Actor) CancelCommand(queue commandPopper, err error) {
	if queue.len() == 0 {
		return
	}
	cmd := queue.pop()
	cmd.replyWithError(err)
	a.activePromises.Delete(cmd.promiseId)
}

func (a *Actor) CancelCommands(queue commandPopper, err error) {
	for queue.len() > 0 {
		cmd := queue.pop()
		cmd.replyWithError(err)
		a.activePromises.Delete(cmd.promiseId)
	}
}

func (a *Actor) Sender() ActorService {
	return a.currentCommand.origin
}

func (a *Actor) GetStatus() Status {
	return Status{len(a.commandProcessors.processors), len(a.messageProcessors.processors),
		len(a.activePromises), len(a.inflightRequests), len(a.streamInputs), len(a.streamOutputs)}
}

//treat ActorCommands as read only, changes to it will only confuse GetInfo requestors but don't change actual commands or handlers
func (a *Actor) GetCommandInfo() ActorCommands {
	return a.commandProcessors.commands
}

func (a *Actor) processCommand(cmd commandMessage) (err errors.StackTraceError) {
	defer func() {
		err = errors.RecoverToError(recover())
		if err != nil {
			a.replyWithError(err)
		}
	}()
	a.currentCommand = cmd
	if _, ok := cmd.data.(GetInfo); ok {
		a.reply(&a.commandProcessors.commands)
		return nil
	}
	response, processErr := a.runCommandProcessor(cmd.data, len(a.commandProcessors.commandFilters))
	if processErr != nil {
		a.replyWithError(processErr)
		return nil
	}
	if response == nil {
		a.reply(nil)
		return nil
	}
	response.Visit((*actorResponseVisitor)(a))
	return nil
}

func (a *Actor) processReissuedCommand(cmd resumedCommand) (err errors.StackTraceError) {
	defer func() {
		err := errors.RecoverToError(recover())
		if err != nil {
			a.replyWithError(err)
		}
	}()
	a.currentCommand = cmd.commandMessage
	a.activePromises.Delete(cmd.promiseId)
	response, processError := a.runCommandProcessor(cmd.data, cmd.filterIndex)
	if processError != nil {
		a.replyWithError(processError)
		return nil
	}
	if response == nil {
		a.reply(nil)
		return nil
	}
	response.Visit((*actorResponseVisitor)(a))
	return nil
}

func (a *Actor) processReissuedCommands() {
	for a.reissuedCommands.len() > 0 {
		cmd := a.reissuedCommands.pop()
		if !cmd.isValid() {
			continue
		}
		err := a.processReissuedCommand(cmd)
		if err != nil {
			a.onPanic(err)
		}
	}
}

func (a *Actor) processServiceFinished(message notifyClose) {
	a.monitoringActors.CallAndRemove(message.destination, message.err)
	a.awaitingActors.CallAndRemove(message.destination, message.err)
}

func (a *Actor) processReply(r reply) {
	if err, ok := r.data.(errorReply); ok {
		a.inflightRequests.Error(r.id, err.err)
	} else {
		a.inflightRequests.Process(r.id, r.data)
	}
}

func (a *Actor) replyFromPromise(id promiseId, data interface{}) {
	if a.activePromises.Contains(id) {
		id.reply(data)
		a.activePromises.Delete(id)
	}
}

func (a *Actor) replyWithErrorFromPromise(id promiseId, err error) {
	a.replyFromPromise(id, errorReply{err})
}

func (a *Actor) cancelCommandProcessingFromPromise(id promiseId, err error) {
	cancel, ok := a.activePromises[id]
	if ok {
		cancel.Call()
		id.replyWithError(err)
		a.activePromises.Delete(id)
	}
}

func (a *Actor) processCancelCommand(command cancelCommand) {
	id := command.toPromiseId()
	cancel, ok := a.activePromises[id]
	if ok {
		cancel.Call()
		a.activePromises.Delete(id)
	}
}

func (a *Actor) processPreReply(message preReply) {
	a.inflightRequests.SetDestination(message.id, message.processor)
}

func (a *Actor) InitStreamRequest(request RequestStream, input StreamInput) {
	a.currentStreamId.Increment()
	base := input.getBase()
	base.init(a, a.currentStreamId)
	input.RequestNext()
	request.setStreamRequest(base.request)
	a.streamInputs.Add(a.currentStreamId, input)
}

//initialize request from input with InitStreamRequest and send it
func (a *Actor) RequestStream(input StreamInput, destination ActorService, request RequestStream, onError common.ErrorCallback) interfaces.Canceller {
	a.InitStreamRequest(request, input)
	return a.SendRequest(destination, request, OnReplyError(func(err error) {
		onError.Call(err)
		a.closeInput(input, err)
	}))
}

func (a *Actor) closeInput(input StreamInput, err error) {
	base := input.getBase()
	enqueue(base.source, downstreamStopped{outputId{base.id, a.Service()}, err})
	a.streamInputs.Remove(base.id)
	base.closed(err)
}

func (a *Actor) processData(id sourceId, data inspect.Inspectable) (err errors.StackTraceError) {
	input := a.streamInputs[id.streamId]
	if input == nil {
		return nil
	}
	defer func() {
		err = errors.RecoverToError(recover())
		if err != nil {
			a.closeInput(input, err)
		}
	}()
	base := input.getBase()
	base.setSource(id.source)
	base.dataReceived()
	processErr := input.Process(data)
	if processErr != nil {
		a.closeInput(input, processErr)
		return nil
	}
	input.RequestNext()
	return nil
}

func (a *Actor) processStreamCanSend(message streamCanSend) errors.StackTraceError {
	return a.processData(sourceId(message), nil)
}

func (a *Actor) processStreamReply(message streamReply) errors.StackTraceError {
	return a.processData(message.id, message.data)
}

func (a *Actor) processUpstreamStopped(message upstreamStopped) {
	input := a.streamInputs[message.id.streamId]
	if input != nil {
		a.streamInputs.Remove(message.id.streamId)
		input.getBase().closed(message.err)
	}
}

func (a *Actor) InitStreamOutput(output StreamOutput, request RequestStream) {
	req := request.getStreamRequest()
	base := output.getBase()
	base.init(a, req.id)
	if req.id.IsValid() {
		if req.data == nil {
			enqueue(req.id.destination, streamCanSend{req.id.streamId, a.Service()})
		} else {
			data, err := output.FillData(req.data, req.maxLen)
			if err != nil {
				a.closeOutput(output, err)
				return
			}
			if data == nil {
				if base.shouldCloseWhenActorCloses && a.state != ActorRunning {
					base.CloseStream(a.quitError)
				}
				if base.isStreamClosing {
					a.closeOutput(output, nil)
					return
				}
				a.streamOutputs.Add(req.id, output)
				out := a.streamOutputs[req.id]
				out.dataRequest = req
				a.streamOutputs[req.id] = out
				return
			}
			enqueue(req.id.destination, streamReply{sourceId{req.id.streamId, a.Service()}, data})
		}
		a.streamOutputs.Add(req.id, output)
	} else {
		output.getBase().streamClosed(ErrBadStream)
	}
}

//send request and expect RequestStream compatible reply.
//If everything is good then call InitStreamOutputFromRequest on reply and output thus initializing stream output
func (a *Actor) RequestStreamOutput(output StreamOutput, destination ActorService, request inspect.Inspectable, replyProcessor interfaces.ReplyProcessor) {
	a.SendRequest(destination, request, OnReply(func(reply interface{}) {
		if streamReply, ok := reply.(RequestStream); ok {
			a.InitStreamOutput(output, streamReply)
		} else {
			output.getBase().streamClosed(ErrNotStreamReply)
		}
		if replyProcessor != nil {
			replyProcessor.Process(reply)
		}
	}).OnError(func(err error) {
		output.getBase().streamClosed(err)
		if replyProcessor != nil {
			replyProcessor.Error(err)
		}
	}))
}

func (a *Actor) processStreamRequest(request streamRequest) (err errors.StackTraceError) {
	out, ok := a.streamOutputs[request.id]
	if !ok {
		return nil
	}
	defer func() {
		err = errors.RecoverToError(recover())
		if err != nil {
			a.closeOutput(out.output, err)
		}
	}()
	out.output.Acknowledged()
	base := out.output.getBase()
	base.acknowledged()
	data, fillErr := out.output.FillData(request.data, request.maxLen)
	if fillErr != nil {
		a.closeOutput(out.output, fillErr)
		return nil
	}
	if data == nil {
		if base.shouldCloseWhenActorCloses && a.state != ActorRunning {
			base.CloseStream(a.quitError)
		}
		if base.isStreamClosing {
			a.closeOutput(out.output, nil)
			return nil
		}
		out.dataRequest = request
		a.streamOutputs[request.id] = out
		return nil
	}
	enqueue(request.id.destination, streamReply{sourceId{request.id.streamId, a.Service()}, data})
	return nil
}

func (a *Actor) processCloseStream(request closeStream) {
	out, ok := a.streamOutputs[outputId(request)]
	if ok {
		out.output.getBase().CloseStream(nil)
	}
}

func (a *Actor) processStreamAcknowledged(request streamAck) (err errors.StackTraceError) {
	out, ok := a.streamOutputs[request.id]
	if !ok {
		return
	}
	defer func() {
		err = errors.RecoverToError(recover())
		if err != nil {
			a.closeOutput(out.output, err)
		}
	}()
	out.output.Acknowledged()
	base := out.output.getBase()
	base.acknowledged()

	if base.isStreamClosing {
		a.closeOutput(out.output, nil)
	}
	return nil
}

func (a *Actor) processDownstreamStopped(request downstreamStopped) {
	out, ok := a.streamOutputs[request.id]
	if !ok {
		return
	}
	a.streamOutputs.Remove(out.output.getBase().outStreamId)
	out.output.getBase().streamClosed(request.err)
}

func (a *Actor) closeOutput(out StreamOutput, err error) {
	base := out.getBase()
	if err != nil {
		base.closeStreamNow(err)
	}
	enqueue(base.outStreamId.destination, upstreamStopped{sourceId{base.outStreamId.streamId, a.Service()}, base.closeError})
	a.streamOutputs.Remove(base.outStreamId)
	base.streamClosed(base.closeError)
	//	fmt.Println("stream closed")
}

func (a *Actor) markOutputReady(id outputId) {
	info, ok := a.streamOutputs[id]
	if ok && info.dataRequest.id.IsValid() {
		a.readyOutputs.Add(info.output)
	}
}

func (a *Actor) flushReadyOutput(info streamOutInfo) (err errors.StackTraceError) {
	defer func() {
		err = errors.RecoverToError(recover())
	}()
	data, fillErr := info.output.FillData(info.dataRequest.data, info.dataRequest.maxLen)
	if fillErr != nil {
		a.closeOutput(info.output, fillErr)
		return nil
	}
	if data != nil {
		enqueue(info.dataRequest.id.destination, streamReply{sourceId{info.dataRequest.id.streamId, a.Service()}, data})
	} else if info.output.getBase().isStreamClosing {
		a.closeOutput(info.output, nil)
	}
	return nil
}

func (a *Actor) quitIfInactive() {
	if a.state == ActorClosed {
		return
	}
	if a.state == ActorRunning && a.isInactive() {
		a.Quit(nil)
	}
	if a.state == ActorQuitting {
		a.runExitProcessor()
		if a.canQuit() {
			a.setState(ActorClosed)
		}
	}
}

func (a *Actor) FlushReadyOutputs() bool {
	for out := range a.readyOutputs {
		info, ok := a.streamOutputs[out.getBase().outStreamId]
		if ok {
			if info.dataRequest.id.IsValid() {
				err := a.flushReadyOutput(info)
				if err != nil {
					a.closeOutput(info.output, err)
					a.onPanic(err)
				}
			}
		}
	}
	a.readyOutputs.Clear()
	a.quitIfInactive()
	return a.state != ActorClosed
}

func (a *Actor) isInactive() bool {
	return a.canQuit() && a.monitoringActors.IsEmpty()
}

func (a *Actor) canQuit() bool {
	//fmt.Printf("%p, have active processors:%v, in: %d, out: %d, monitor: %d\n", a.Service(), a.haveActiveProcessors(), len(a.streamInputs), len(a.streamOutputs), len(a.monitoringActors))
	//	fmt.Printf("%p,active promises are empty:%v\n", a.Service(),a.activePromises.IsEmpty())
	//	fmt.Printf("%p,inflight requests are empty:%v\n", a.Service(),a.inflightRequests.IsEmpty())
	return !a.haveActiveProcessors() && a.activePromises.IsEmpty() && a.inflightRequests.IsEmpty() &&
		a.awaitingActors.IsEmpty() && a.streamInputs.IsEmpty() && a.streamOutputs.IsEmpty()
}

func (a *Actor) processIncomingMessage(msg interface{}) (err errors.StackTraceError) {
	debug.Printf("- %p [%s]: %#v\n", a.Service(), a.GetCommandInfo().Name, msg)
	defer func() {
		err2 := errors.RecoverToError(recover())
		if err2 != nil {
			err = err2
		}
	}()
	switch message := msg.(type) {
	case commandMessage:
		err = a.processCommand(message)
	case notifyClose:
		a.processServiceFinished(message)
	case reply:
		a.processReply(message)
	case preReply:
		a.processPreReply(message)
	case quitMessage:
		a.Quit(message.err)
	case closeMessage:
		a.setState(ActorClosed)
	case cancelCommand:
		a.processCancelCommand(message)
	case establishLink:
		a.links.Add(message.source, message.linkType)
	case streamCanSend:
		err = a.processStreamCanSend(message)
	case streamRequest:
		err = a.processStreamRequest(message)
	case streamReply:
		err = a.processStreamReply(message)
	case streamAck:
		err = a.processStreamAcknowledged(message)
	case upstreamStopped:
		a.processUpstreamStopped(message)
	case downstreamStopped:
		a.processDownstreamStopped(message)
	case closeStream:
		a.processCloseStream(message)
	case subscribeStateChange:
		a.stateChangeSubscribers.Add(message.source)
	case notifyStateChange:
		a.runStateChangeProcessor(message.source, message.state)
	default:
		a.runMessageProcessor(msg)
	}
	a.FlushReadyOutputs()
	a.processReissuedCommands()
	return err
}

func (a *Actor) onPanic(err errors.StackTraceError) {
	if a.panicProcessor != nil {
		a.panicProcessor(err)
	} else {
		a.Quit(err)
	}
}

func (a *Actor) readIncomingMessages(head *queue.QueueElement) {
	if head == nil {
		return
	}
	for ; head != nil; head = head.Next {
		if array, ok := head.Data.([]interface{}); ok {
			for i := len(array) - 1; i >= 0; i-- {
				a.messages = append(a.messages, array[i])
			}
		} else if array, ok := head.Data.([]inspect.Inspectable); ok {
			for i := len(array) - 1; i >= 0; i-- {
				a.messages = append(a.messages, array[i])
			}
		} else {
			a.messages = append(a.messages, head.Data)
		}
	}
}

func (a *Actor) processIncomingMessages(head *queue.QueueElement) bool {
	a.readIncomingMessages(head)
	for i := len(a.messages) - 1; i >= 0; i-- {
		err := a.processIncomingMessage(a.messages[i])
		if err != nil {
			a.onPanic(err)
		}
	}
	a.messages = a.messages[:0]
	a.quitIfInactive()
	return a.state != ActorClosed
}

func (a *Actor) IncomingChannel() common.OutSignalChannel {
	return a.queue.SignalChannel()
}

func (a *Actor) ProcessMessages() bool {
	return a.processIncomingMessages(a.queue.TakeHead())
}

func (a *Actor) Run() error {
	for _ = range a.IncomingChannel() {
		if !a.ProcessMessages() {
			return nil
		}
	}
	return nil
}

func (a *Actor) waitDependentActors() {
	debug.Printf("%p[z], %d links remaining\n", a.Service(), len(a.links))
	for _ = range a.IncomingChannel() {
		a.readIncomingMessages(a.queue.TakeHead())
		for i := len(a.messages) - 1; i >= 0; i-- {
			msg := a.messages[i]
			debug.Printf("- %p[z]: %#v\n", a.Service(), msg)
			if notify, ok := msg.(notifyClose); ok {
				a.links.Remove(notify.destination)
				debug.Printf("%p[z], %d links remaining\n", a.Service(), len(a.links))
				continue
			}
			if subscribe, ok := msg.(subscribeStateChange); ok {
				a.stateChangeSubscribers.Add(subscribe.source)
				continue
			}
			zombieBehaviour(a.Service(), msg, ErrActorDead)
		}
		a.messages = a.messages[:0]
		if a.links.IsEmpty() {
			break
		}
	}
}
