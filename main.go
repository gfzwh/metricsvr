package main

import (
	"metricsvr/controller"
	"os"
	"os/signal"
	"syscall"

	"github.com/shockerjue/gffg/config"
	"github.com/shockerjue/gffg/zzlog"
)

func main() {
	config.Init("./conf/gffg.xml")
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
