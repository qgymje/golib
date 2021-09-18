package utility

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForTerminal(fn func()) {
	stopSignals := make(chan os.Signal, 1)
	signal.Notify(stopSignals, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignals

	fn()
}
