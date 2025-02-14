package plugins

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Plugin interface {
	Name() string
	Commands() []string
	Handle(client *whatsmeow.Client, msg *events.Message)
}

var registeredPlugins []Plugin

func Register(plugin Plugin) {
	registeredPlugins = append(registeredPlugins, plugin)
}

func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
	if msg.Info.IsFromMe || msg.Message.GetConversation() == "" {
		return
	}

	text := msg.Message.GetConversation()
	for _, plugin := range registeredPlugins {
		for _, cmd := range plugin.Commands() {
			if text == cmd {
				plugin.Handle(client, msg)
				return
			}
		}
	}
}
