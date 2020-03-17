package flags

import (
	"net/url"
	"strconv"
	"strings"
)

type HostPort struct {
	Host string
	Port int
}

func (hp *HostPort) SetDefaultsIfEmpty(defaultHost string, defaultPort int) {
	if len(hp.Host) == 0 {
		hp.Host = defaultHost
	}
	if hp.Port == 0 {
		hp.Port = defaultPort
	}
}

func (hp *HostPort) LoadFromUrl(url *url.URL) {
	portString := url.Port()
	if len(portString) == 0 {
		hp.Port = 0
	} else {
		port, err := strconv.ParseInt(portString, 10, 32)
		if err != nil {
			hp.Port = 0
		} else {
			hp.Port = int(port)
		}
	}
	hp.Host = url.Hostname()
}

func (hp *HostPort) LoadFromString(host string) {
	split := strings.Split(host, ":")
	switch len(split) {
	case 1:
		hp.Host = split[0]
		hp.Port = 0
	case 2:
		port, err := strconv.ParseInt(split[1], 10, 32)
		if err != nil {
			hp.Host = ""
			hp.Port = 0
		} else {
			hp.Host = split[0]
			hp.Port = int(port)
		}
	default:
		hp.Host = ""
		hp.Port = 0
	}
}

func (hp *HostPort) RegisterFlags(prefix string) {
	if len(prefix) > 0 {
		StringFlag(&hp.Host, prefix+"-host", prefix+" host")
		IntFlag(&hp.Port, prefix+"-port", prefix+" port")
	} else {
		StringFlag(&hp.Host, "host", "host")
		IntFlag(&hp.Port, "port", "port")
	}
}

func (hp *HostPort) RegisterFlagsWithDescriptions(prefix string, hostDescription string, portDescription string) {
	if len(prefix) > 0 {
		StringFlag(&hp.Host, prefix+"-host", hostDescription)
		IntFlag(&hp.Port, prefix+"-port", portDescription)
	} else {
		StringFlag(&hp.Host, "host", hostDescription)
		IntFlag(&hp.Port, "port", portDescription)
	}
}

func (hp HostPort) String() string {
	return hp.Host + ":" + strconv.FormatInt(int64(hp.Port), 10)
}

func (hp HostPort) StringArray() []string {
	return []string{hp.String()}
}
