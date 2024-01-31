package main

import (
	"metricsnode/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfzwh/gfz/udp"
)

func main() {
	svr := udp.NewServer("./conf/gfz.xml")
	defer svr.Release()

	ctl := controller.Controller()
	svr.Run(udp.Event(ctl.OnEvent))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	return
}
