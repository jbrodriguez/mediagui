package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
)

type Core struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	// socket   *Socket

	mailbox chan *pubsub.Mailbox
}

func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{bus: bus, settings: settings}
	core.init()
	return core
}

func (c *Core) Start() {
	mlog.Info("Starting service Core ...")

	c.mailbox = c.register(c.bus, c.getConfig, "/get/config")
	// c.mailbox = c.bus.Sub("/get/config")

	go c.react()
}

func (c *Core) Stop() {
	mlog.Info("Stopped service Core ...")
}

func (c *Core) react() {
	for mbox := range c.mailbox {
		mlog.Info("Core:Topic: %s", mbox.Topic)
		c.dispatch(mbox.Topic, mbox.Content)
	}
}

func (c *Core) getConfig(msg *pubsub.Message) {
	msg.Reply <- &c.settings.Config
	mlog.Info("Sent config")
}
