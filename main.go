package main

import (
	"metricsvr/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfzwh/gfz/config"
)

func main() {
	config.Init("./conf/gfz.xml")
	ctl := controller.Controller()
	go ctl.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	return
}
