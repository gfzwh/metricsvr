package controller

import (
	"github.com/gfzwh/gfz/config"
	"github.com/gfzwh/gfz/kafka"
)

type controller struct {
	sub *kafka.Consumer
}

func newController() *controller {
	v := &controller{
		sub: kafka.NewConsumer(
			kafka.Brokers(config.Get("metrics", "brokers").String("")),
			kafka.Group(config.Get("metrics", "group").String("")),
			kafka.Topic(config.Get("metrics", "topic").String(""))),
	}

	return v
}

func Controller() *controller {
	return newController()
}

func (this *controller) Release() {
	this.sub.Release()
}

func (this *controller) Run() {
	this.sub.Consume(this)
}
