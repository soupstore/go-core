package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/soupstore/go-core/logging"
)

// GracefulShutdownOnSignal will run the supplied shutdown procedure when one of the specified signals is received
func GracefulShutdownOnSignal(signals []syscall.Signal, shutdownProcedure func()) {
	stopChan := make(chan os.Signal)
	for _, s := range signals {
		signal.Notify(stopChan, s)
	}

	<-stopChan
	logging.Info("Shutting down...")

	shutdownProcedure()

	logging.Info("Shut down complete")
}
