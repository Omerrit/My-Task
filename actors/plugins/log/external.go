package log

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
)

type WriterId int

var writerCounter WriterId

func RegisterWriter() (id WriterId) {
	writerCounter++
	return writerCounter - 1
}

var verbosityFlag int

func GetVerbosity() Verbosity {
	return Verbosity(verbosityFlag)
}

func SubscribeForMessages(actor actors.ActorCompatible, id WriterId, input actors.StreamInput) {
	base := actor.GetBase()
	base.RequestStream(input, base.System().GetPluginActor(pluginName), &subscribe{id: int(id)}, base.Quit)
}

func NotifyNotSubscribe(actor actors.ActorCompatible, id WriterId) {
	base := actor.GetBase()
	base.SendMessage(base.System().GetPluginActor(pluginName), &notGonnaSubscribe{id: int(id)})
}

func Debugln(actor actors.ActorCompatible, values ...interface{}) {
	writeSimpleMessage(actor, SeverityDebug, errors.CallerStack(defaultStackDepth), fmt.Sprintln(values...))
}

func Debugf(actor actors.ActorCompatible, format string, values ...interface{}) {
	writeSimpleMessage(actor, SeverityDebug, errors.CallerStack(defaultStackDepth), fmt.Sprintf(format, values...))
}

func DebugFull(actor actors.ActorCompatible, source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(actor, SeverityDebug, errors.CallerStack(defaultStackDepth), source, message, extra, sourceInfo)
}

func DebugErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityDebug, errors.CallerStack(defaultStackDepth))
}

func Infoln(actor actors.ActorCompatible, values ...interface{}) {
	writeSimpleMessage(actor, SeverityInfo, errors.CallerStack(defaultStackDepth), fmt.Sprintln(values...))
}

func Infof(actor actors.ActorCompatible, format string, values ...interface{}) {
	writeSimpleMessage(actor, SeverityInfo, errors.CallerStack(defaultStackDepth), fmt.Sprintf(format, values...))
}

func InfoFull(actor actors.ActorCompatible, source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(actor, SeverityInfo, errors.CallerStack(defaultStackDepth), source, message, extra, sourceInfo)
}

func InfoErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityInfo, errors.CallerStack(defaultStackDepth))
}

func Statusln(actor actors.ActorCompatible, values ...interface{}) {
	writeSimpleMessage(actor, SeverityStatus, errors.CallerStack(defaultStackDepth), fmt.Sprintln(values...))
}

func Statusf(actor actors.ActorCompatible, format string, values ...interface{}) {
	writeSimpleMessage(actor, SeverityStatus, errors.CallerStack(defaultStackDepth), fmt.Sprintf(format, values...))
}

func StatusFull(actor actors.ActorCompatible, source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(actor, SeverityStatus, errors.CallerStack(defaultStackDepth), source, message, extra, sourceInfo)
}

func StatusErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityStatus, errors.CallerStack(defaultStackDepth))
}

func Processingln(actor actors.ActorCompatible, sourceInfo interface{}, values ...interface{}) {
	writeMessage(actor, SeverityProcessing, errors.CallerStack(defaultStackDepth), "", fmt.Sprintln(values...), nil, sourceInfo)
}

func Processingf(actor actors.ActorCompatible, sourceInfo interface{}, format string, values ...interface{}) {
	writeMessage(actor, SeverityProcessing, errors.CallerStack(defaultStackDepth), "", fmt.Sprintf(format, values...), nil, sourceInfo)
}

func ProcessingFull(actor actors.ActorCompatible, sourceInfo interface{}, source string, message string, extra interface{}) {
	writeMessage(actor, SeverityProcessing, errors.CallerStack(defaultStackDepth), source, message, extra, sourceInfo)
}

func ProcessingErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityProcessing, errors.CallerStack(defaultStackDepth))
}

func Warningln(actor actors.ActorCompatible, values ...interface{}) {
	writeSimpleMessage(actor, SeverityWarning, errors.CallerStack(defaultStackDepth), fmt.Sprintln(values...))
}

func Warningf(actor actors.ActorCompatible, format string, values ...interface{}) {
	writeSimpleMessage(actor, SeverityWarning, errors.CallerStack(defaultStackDepth), fmt.Sprintf(format, values...))
}

func WarningFull(actor actors.ActorCompatible, source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(actor, SeverityWarning, errors.CallerStack(defaultStackDepth), source, message, extra, sourceInfo)
}

func WarningErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityWarning, errors.CallerStack(defaultStackDepth))
}

func ErrorErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityError, errors.CallerStack(defaultStackDepth))
}

func CriticalErr(actor actors.ActorCompatible, err error) {
	writeErrorNoInfo(actor, err, SeverityCritical, errors.CallerStack(defaultStackDepth))
}
