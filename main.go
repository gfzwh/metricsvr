package main

import (
	"metricsvr/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfzwh/gfz/config"
	"github.com/gfzwh/gfz/zzlog"
)

func main() {
	config.Init("./conf/gfz.xml")
	zzlog.Init(
		zzlog.WithLogName(config.Get("log", "log_file").String("")),
		zzlog.WithLevel(config.Get("log", "level").String("info")))

	ctl := controller.Controller()
	defer ctl.Release()

	go ctl.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	return
}
