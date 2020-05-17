package starter

import (
	"flag"
	"fmt"
	"gerrit-share.lan/go/actors"
	"gerrit-share.lan/go/utils/flags"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"time"
)

const defaultAutorestartPeriod = 5 * time.Second
const defaultLogFile = ""

func runInitializers() {
	var names []string
	for name, initializer := range serviceFlagInitializers {
		if initializer != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		serviceFlagInitializers[name]()
	}
}

func logAlsoToFile(filename string) (*os.File, error) {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_SYNC|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	return logFile, nil
}

func Launch(id string) {
	if len(id) == 0 {
		id = filepath.Base(os.Args[0])
	}
	flags.StringFlag(&id, "id", "executable id")
	autorestartPeriodSec := int(defaultAutorestartPeriod / time.Second)
	flags.IntFlag(&autorestartPeriodSec, "autorestart-period", "autorestart period for crashed always on services, seconds")
	logFile := defaultLogFile
	flags.StringFlag(&logFile, "log", "also write log to the specified file in addition to console")
	enableSurveillance := false
	flags.BoolFlag(&enableSurveillance, "enable-surveillance", "run surveillance actor for more output")
	runInitializers()
	flag.Parse()
	for _, processor := range serviceFlagProcessors {
		processor()
	}
	if len(logFile) > 0 {
		file, err := logAlsoToFile(logFile)
		if err != nil {
			fmt.Println("Failed to open log file", err)
			return
		}
		defer file.Close()
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	var system actors.System
	if enableSurveillance {
		log.Println("Enabling surveillance")
		system.EnableSurveillance()
	}
	service, err := newStarterService(&system, starterServiceName, time.Duration(autorestartPeriodSec)*time.Second)
	if err != nil {
		log.Println("failed to start starter service", err)
		system.WaitFinished()
	}
	fmt.Println("Running, press ctrl-c to exit")
	select {
	case <-service.DoneChannel():
		log.Println("Main service quit,", service.CloseError())
		system.WaitFinished()
		return
	case <-interrupt:
	}

	log.Println("Shutting down, press ctrl-c to to quit without properly shutting down (you may lose data, logs and hell knows what else)")
	service.SendQuit(nil)
	go func() {
		<-interrupt
		log.Println("You've been warned")
		os.Exit(1)
	}()
	system.WaitFinished()
	log.Println("Normal shutdown done", service.CloseError())
}
