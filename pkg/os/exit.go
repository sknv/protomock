package os

import (
	"os"
	"os/signal"
	"syscall"
)

// NotifyAboutExit returns a channel to catch a program exit signal.
func NotifyAboutExit() <-chan os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(
		exit,
		os.Interrupt, syscall.SIGTERM,
	)

	return exit
}
