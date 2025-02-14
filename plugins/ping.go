package plugins

import (
    "context"
    "fmt"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "go.mau.fi/whatsmeow/types/events"
    "google.golang.org/protobuf/proto"
)

// init registers the "ping" and "!ping" commands.
func init() {
    RegisterCommand("ping", PingHandler)
    RegisterCommand("!ping", PingHandler)
}

// PingHandler responds to a ping command by sending a "pong" message.
func PingHandler(client *whatsmeow.Client, evt *events.Message) {
    // Build a text message using waE2E.Message.
    msg := &waE2E.Message{
	Conversation: proto.String("pong"),
    }

    // Determine the target chat.
    target := evt.Info.Chat
    if target.IsEmpty() {
	target = evt.Info.Sender
    }

    resp, err := client.SendMessage(context.Background(), target, msg)
    if err != nil {
	fmt.Println("Error sending pong response:", err)
	return
    }
    fmt.Println("Pong sent at:", resp.Timestamp)
}