package gd

import (
	"os"
	"os/signal"
	"syscall"
)

// ShutDownSignals returns all the signals that are being watched for to shut down services.
func ShutDownSignals() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
	}
}

func ListenShutDownSignals(signals ...os.Signal) {
	ch := make(chan os.Signal)
	signal.Notify(ch, signals...)
	for {
		select {
		case <-ch:
			os.Exit(0)
		}
	}
}
