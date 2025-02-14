package lib

import (
    "fmt"
    "time"

    "go.mau.fi/whatsmeow/types/events"
)

// PrintIncomingMessage logs details of an incoming message.
func PrintIncomingMessage(evt *events.Message) {
    // Determine sender.
    sender := evt.Info.Sender.String()

    // Determine chat type: if Chat is empty or equals sender, it's private.
    chatJID := evt.Info.Chat
    chatType := "private"
    if !chatJID.IsEmpty() && chatJID.String() != sender {
	chatType = "group"
    }
    chat := chatJID.String()
    if chat == "" {
	chat = sender // fallback if Chat is empty
    }

    // Use the event's timestamp if available; otherwise, use current time.
    ts := evt.Info.Timestamp
    if ts.IsZero() {
	ts = time.Now()
    }

    // Get the text content.
    content := evt.Message.GetConversation()

    // Print a formatted log.
    fmt.Printf("[INCOMING] From: %s | Chat: %s (%s) | At: %s\nMessage: %s\n", sender, chat, chatType, ts.Format(time.RFC1123), content)
}
