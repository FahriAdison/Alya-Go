package plugins

import (
    "context"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// Ping plugin
type PingPlugin struct{}

func (p *PingPlugin) Name() string {
    return "Ping"
}

func (p *PingPlugin) Commands() []string {
    return []string{"!ping"}
}

func (p *PingPlugin) Handle(client *whatsmeow.Client, msg *events.Message) {
    _, err := client.SendMessage(context.Background(), msg.Info.Chat, &whatsmeow.TextMessage{
	Content: "🏓 Pong!",
    })
    if err != nil {
	println("Failed to send ping response:", err.Error())
    }
}

// Register plugin
func init() {
    Register(&PingPlugin{})
}