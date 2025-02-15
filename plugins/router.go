package plugins

import (
    "strings"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

type CommandHandler func(client *whatsmeow.Client, evt *events.Message)

var commandRegistry = make(map[string]CommandHandler)

// RegisterCommand adds a command to the registry
func RegisterCommand(cmd string, handler CommandHandler) {
    commandRegistry[strings.ToLower(cmd)] = handler
}

// Handle processes incoming messages and routes commands
func Handle(client *whatsmeow.Client, evt *events.Message) {
    text := evt.Message.GetConversation()

    // ✅ Check if it's an image caption
    if text == "" && evt.Message.GetImageMessage() != nil {
	text = evt.Message.GetImageMessage().GetCaption()
    }

    // Ignore empty messages
    if text == "" {
	return
    }

    // Extract the command
    words := strings.Fields(text)
    command := strings.ToLower(words[0])

    // Run command if found
    if handler, ok := commandRegistry[command]; ok {
	handler(client, evt)
    }
}
