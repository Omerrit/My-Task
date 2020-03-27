package log

import (
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/errors"
)

type Logger struct {
	actors.Actor
	sourceInfo interface{}
	source     string
	stackDepth uint
}

func (l *Logger) SetLogSource(source string) {
	l.source = source
}

func (l *Logger) SetLogSourceInfo(info interface{}) {
	l.sourceInfo = info
}

func (l *Logger) SetLogStackDepth(depth int) {
	l.stackDepth = uint(depth)
}

func (l *Logger) Debugln(values ...interface{}) {
	writeMessage(l, SeverityDebug, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintln(values...), nil, l.sourceInfo)
}

func (l *Logger) Debugf(format string, values ...interface{}) {
	writeMessage(l, SeverityDebug, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintf(format, values...), nil, l.sourceInfo)
}

func (l *Logger) DebugFull(source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(l, SeverityDebug, errors.CallerStack(l.stackDepth), source, message, extra, sourceInfo)
}

func (l *Logger) DebugErr(err error) {
	writeErrorWrapper(l, err, SeverityDebug, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) Infoln(values ...interface{}) {
	writeMessage(l, SeverityInfo, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintln(values...), nil, l.sourceInfo)
}

func (l *Logger) Infof(format string, values ...interface{}) {
	writeMessage(l, SeverityInfo, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintf(format, values...), nil, l.sourceInfo)
}

func (l *Logger) InfoFull(source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(l, SeverityInfo, errors.CallerStack(l.stackDepth), source, message, extra, sourceInfo)
}

func (l *Logger) InfoErr(err error) {
	writeErrorWrapper(l, err, SeverityInfo, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) Statusln(values ...interface{}) {
	writeMessage(l, SeverityStatus, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintln(values...), nil, l.sourceInfo)
}

func (l *Logger) Statusf(format string, values ...interface{}) {
	writeMessage(l, SeverityStatus, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintf(format, values...), nil, l.sourceInfo)
}

func (l *Logger) StatusFull(source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(l, SeverityStatus, errors.CallerStack(l.stackDepth), source, message, extra, sourceInfo)
}

func (l *Logger) StatusErr(err error) {
	writeErrorWrapper(l, err, SeverityStatus, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) Processingln(values ...interface{}) {
	writeMessage(l, SeverityProcessing, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintln(values...), nil, l.sourceInfo)
}

func (l *Logger) Processingf(format string, values ...interface{}) {
	writeMessage(l, SeverityProcessing, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintf(format, values...), nil, l.sourceInfo)
}

func (l *Logger) ProcessingFull(sourceInfo interface{}, source string, message string, extra interface{}) {
	writeMessage(l, SeverityProcessing, errors.CallerStack(l.stackDepth), source, message, extra, sourceInfo)
}

func (l *Logger) ProcessingErr(err error) {
	writeErrorWrapper(l, err, SeverityProcessing, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) Warningln(values ...interface{}) {
	writeMessage(l, SeverityWarning, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintln(values...), nil, l.sourceInfo)
}

func (l *Logger) Warningf(format string, values ...interface{}) {
	writeMessage(l, SeverityWarning, errors.CallerStack(l.stackDepth), l.source, fmt.Sprintf(format, values...), nil, l.sourceInfo)
}

func (l *Logger) WarningFull(source string, message string, extra interface{}, sourceInfo interface{}) {
	writeMessage(l, SeverityWarning, errors.CallerStack(l.stackDepth), source, message, extra, sourceInfo)
}

func (l *Logger) WarningErr(err error) {
	writeErrorWrapper(l, err, SeverityWarning, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) ErrorErr(err error) {
	writeErrorWrapper(l, err, SeverityError, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}

func (l *Logger) CriticalErr(err error) {
	writeErrorWrapper(l, err, SeverityCritical, errors.CallerStack(l.stackDepth), l.source, l.sourceInfo)
}
