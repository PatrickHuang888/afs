package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/PatrickHuang888/afs/daemon"
)

var DaemonCommand = cli.Command{
	Name:   "daemon",
	Usage:  "start a long-running gfs daemon process",
	Action: daemonCommand,
}

func daemonCommand(c *cli.Context) error {
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	notify := make(chan os.Signal, 2)
	signal.Notify(notify, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer signal.Stop(notify)

		<-notify
		cancel()

		select {
		case <-time.After(30 * time.Second):
			fmt.Println("Timed out on shutdown, terminating...")
		case <-notify:
			fmt.Println("Received another interrupt before graceful shutdown, terminating...")
		}
		os.Exit(-1)
	}()

	srv, err := daemon.New(ctx, "127.0.0.1", "2345")
	if err != nil {
		return err
	}
	return srv.Run()
}
