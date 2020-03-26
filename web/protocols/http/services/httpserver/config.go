package httpserver

import (
	"gerrit-share.lan/go/inspect"
	"time"
)

const DefaultHttpServerName = "http_server"

const packageName = "http"

const (
	endpointsFileName          = "/endpoints"
	defaultJwtKey              = `p+}Je2tN&>/8q=GDk{@JQtW|Rc4fiX+R|{P5pZkzUUNEm11og4%\~AJ}ET2Y}Q]`
	defaultSessionDuration     = time.Minute * 525600
	defaultSessionResetPercent = 0.1
)

var defaultConfig = Config{
	JwtKey:              defaultJwtKey,
	SessionDuration:     defaultSessionDuration,
	SessionResetPercent: defaultSessionResetPercent,
}

// httpServer Config
type Config struct {
	// secret key used to sign cookie's value (jwt)
	JwtKey string
	// duration of session (both jwt and session cookie)
	SessionDuration time.Duration
	// remaining session percentage when httpServer automatically refresh session cookie
	SessionResetPercent float64
}

const configName = packageName + ".Config"

func (c *Config) Inspect(i *inspect.GenericInspector) {
	objectInspector := i.Object(configName, "http server Config")
	{
		objectInspector.String(&c.JwtKey, "jwt", false, "jwt secret key")
		objectInspector.Float64(&c.SessionResetPercent, 'g', -1, "sessionreset", false, "session reset percent")
		if objectInspector.IsReading() {
			var str string
			var err error
			objectInspector.String(&str, "session", false, "session duration")
			if len(str) == 0 {
				c.SessionDuration = defaultSessionDuration
			} else {
				c.SessionDuration, err = time.ParseDuration(str)
				if err != nil {
					objectInspector.SetError(err)
					c.SessionDuration = defaultSessionDuration
				}
			}
			if c.SessionResetPercent <= 0 || c.SessionResetPercent >= 1 {
				c.SessionResetPercent = defaultSessionResetPercent
			}
			if len(c.JwtKey) == 0 {
				c.JwtKey = defaultJwtKey
			}
		} else {
			str := c.SessionDuration.String()
			objectInspector.String(&str, "session", false, "session duration")
		}
		objectInspector.End()
	}
}
