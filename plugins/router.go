package plugins

import (
    "strings"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

type CommandHandler func(client *whatsmeow.Client, evt *events.Message)

var commandRegistry = make(map[string]CommandHandler)

func RegisterCommand(cmd string, handler CommandHandler) {
    commandRegistry[strings.ToLower(cmd)] = handler
}

func Handle(client *whatsmeow.Client, evt *events.Message) {
    text := evt.Message.GetConversation()
    if text == "" {
	return
    }
    words := strings.Fields(text)
    command := strings.ToLower(words[0])
    if handler, ok := commandRegistry[command]; ok {
	handler(client, evt)
    }
}
