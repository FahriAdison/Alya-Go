package plugins

import (
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
)

// Plugin interface
type Plugin interface {
    Name() string
    Commands() []string
    Handle(client *whatsmeow.Client, msg *events.Message)
}

// Registered plugins
var registeredPlugins []Plugin

// Register a plugin
func Register(plugin Plugin) {
    registeredPlugins = append(registeredPlugins, plugin)
}

// Handle incoming messages
func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
    if msg.Info.IsFromMe || msg.Message.GetConversation() == "" {
	return
    }

    // Check for commands
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