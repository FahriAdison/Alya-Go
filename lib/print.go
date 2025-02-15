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

    // Get the message content.
    var content string
    if evt.Message.GetConversation() != "" {
        // Text message
        content = evt.Message.GetConversation()
    } else if evt.Message.GetImageMessage() != nil {
        // Image message
        img := evt.Message.GetImageMessage()
        content = fmt.Sprintf("[Image] Caption: %s", img.GetCaption())
    } else if evt.Message.GetVideoMessage() != nil {
        // Video message
        vid := evt.Message.GetVideoMessage()
        content = fmt.Sprintf("[Video] Caption: %s", vid.GetCaption())
    } else if evt.Message.GetDocumentMessage() != nil {
        // Document message
        doc := evt.Message.GetDocumentMessage()
        content = fmt.Sprintf("[Document] Filename: %s", doc.GetFileName())
    } else {
        // Other message types
        content = "[Unsupported Message Type]"
    }

    // Print a formatted log.
    fmt.Printf("[INCOMING] From: %s | Chat: %s (%s) | At: %s\nMessage: %s\n", sender, chat, chatType, ts.Format(time.RFC1123), content)
}