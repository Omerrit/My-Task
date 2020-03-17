package tasks

import (
	"flag"
	"gerrit-share.lan/go/actors/starter"
	"gerrit-share.lan/go/utils/flags"
	"net/url"
)

const serviceName = "task_storage"

var postgresUrl url.URL

func init() {
	var hostPort flags.HostPort
	var user string
	var password string
	var useSsl bool

	starter.SetFlagInitializer(serviceName, func() {
		hostPort.LoadFromUrl(&postgresUrl)
		hostPort.SetDefaultsIfEmpty("localhost", 5432)
		flag.StringVar(&hostPort.Host, "postgres-host", hostPort.Host, "postgresql host")
		flag.IntVar(&hostPort.Port, "postgres-port", hostPort.Port, "postgresql port")
		flag.StringVar(&user, "postgres-user", postgresUrl.User.Username(), "postgresql user")
		password, _ = postgresUrl.User.Password()
		flag.StringVar(&password, "postgres-password", password, "postgresql password")
		if len(postgresUrl.Path) == 0 {
			postgresUrl.Path = "db"
		}
		flag.StringVar(&postgresUrl.Path, "postgres-database", postgresUrl.Path, "postgresql database name")
		flag.BoolVar(&useSsl, "postgres-use-ssl", useSsl, "use ssl when connecting to postgresql server")
	})
	starter.SetFlagProcessor(serviceName, func() {
		postgresUrl.Scheme = "postgresql"
		postgresUrl.Host = hostPort.String()
		if len(user) != 0 {
			if len(password) != 0 {
				postgresUrl.User = url.UserPassword(user, password)
			} else {
				postgresUrl.User = url.User(user)
			}
		} else {
			postgresUrl.User = nil
		}
		values := postgresUrl.Query()
		if useSsl {
			values.Add("sslmode", "require")
		} else {
			values.Add("sslmode", "disable")
		}
		postgresUrl.RawQuery = values.Encode()
	})
}
