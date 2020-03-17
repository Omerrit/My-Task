package services

import (
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/utils/flags"
	"os"
)

const DefaultHttpServerName = "http_server"

var (
	defaultHttpServerParams = flags.HostPort{Port: 8882}
	defaultDir, _           = os.Getwd()
)

func DefaultHttpServerParams() flags.HostPort {
	return defaultHttpServerParams
}

func DefaultDir() string {
	return defaultDir
}

func init() {
	starter.SetCreator(DefaultHttpServerName, nil)
	starter.SetFlagInitializer(DefaultHttpServerName, func() {
		defaultHttpServerParams.RegisterFlagsWithDescriptions(
			"http",
			"listen to http requests on this hostname/ip address",
			"listen to http requests on this port")
		flags.StringFlag(&defaultDir, "endpoints", "saved endpoints directory")
	})
}
