package timetypes

import (
	"gerrit-share.lan/go/inspect"
	"time"
)

const packageName = "time"

type UnixMilliseconds struct {
	time.Time
}

const UnixMillisecondsName = packageName + ".msecs"

func (u *UnixMilliseconds) Inspect(inspector *inspect.GenericInspector) {
	const description = "unix timestamp,milliseconds"
	if inspector.IsReading() {
		var t int64
		inspector.Int64(&t, UnixMillisecondsName, description)
		secs := t / 1e3
		frac := t % 1e3
		u.Time = time.Unix(secs, frac*1e6)
	} else {
		t := u.Unix() + int64(u.Nanosecond()/1e6)
		inspector.Int64(&t, UnixMillisecondsName, description)
	}
}
