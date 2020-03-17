package starter

import (
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/actors/plugins/registry"
	"gerrit-share.lan/go/actors/services/autorestarter"
	"time"
)

func newStarterService(system *actors.System, serviceName string, autorestartPeriod time.Duration) (actors.ActorService, error) {
	var launchers autorestarter.ServiceMakers
	for name, maker := range defaultServiceCreators {
		if maker != nil {
			launchers.Add(name, autorestarter.ServiceMaker(maker))
		}
	}
	result := autorestarter.NewStaticAutorestarter(system, serviceName, launchers, autorestartPeriod)
	err := system.Do(func(actor *actors.Actor) {
		registry.RegisterOther(actor, serviceName, result, func(err error) {
			result.SendQuit(nil)
			actor.Quit(err)
		})
	}).CloseError()
	return result, err
}
