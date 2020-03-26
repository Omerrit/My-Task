package cookies

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

const sessionCookieName = "session"

type CookieInfo struct {
	SessionId       uuid.UUID
	ExpiresAt       int64
	SessionReset    float64
	SessionDuration time.Duration
}

func AddCookie(w http.ResponseWriter, sessionId []byte, secretKey string, duration time.Duration) {
	if sessionId == nil {
		return
	}
	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   sign(sessionId, secretKey, duration),
		Expires: time.Now().Add(duration),
	}
	http.SetCookie(w, &cookie)
}

func ParseCookie(request *http.Request, secretKey string) ([]byte, int64, error) {
	c, err := request.Cookie(sessionCookieName)
	if err != nil {
		return nil, 0, nil
	}
	return parse(c.Value, secretKey)
}
