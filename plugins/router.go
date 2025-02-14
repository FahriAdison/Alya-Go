package plugins

import (
    "strings"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// CommandHandler defines the function signature for command handlers.
// Now it takes a *whatsmeow.Client instead of an interface{}.
type CommandHandler func(client *whatsmeow.Client, evt *events.Message)

// commandRegistry maps command keywords (lowercased) to their handler functions.
var commandRegistry = make(map[string]CommandHandler)

// RegisterCommand registers a command with its handler.
func RegisterCommand(cmd string, handler CommandHandler) {
    commandRegistry[strings.ToLower(cmd)] = handler
}

// Handle dispatches an incoming message to a registered command handler if the first word matches.
func Handle(client *whatsmeow.Client, evt *events.Message) {
    text := evt.Message.GetConversation()
    if text == "" {
	return
    }

    words := strings.Fields(text)
    command := strings.ToLower(words[0])
    if handler, exists := commandRegistry[command]; exists {
	handler(client, evt)
    }
}