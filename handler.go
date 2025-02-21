package main

import (
    "fmt"
    "log"
    "strings"
    "sync"
    "time"

    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
)

type MessageHandler struct {
    client  *whatsmeow.Client
    mutex   sync.RWMutex
    cooldowns map[string]int64
}

func NewMessageHandler(client *whatsmeow.Client) *MessageHandler {
    return &MessageHandler{
        client:    client,
        cooldowns: make(map[string]int64),
    }
}

// HandleMessage processes incoming messages
func (h *MessageHandler) HandleMessage(evt *events.Message) {
    // Ignore messages from self
    if evt.Info.IsFromMe {
        return
    }

    // Get sender info
    sender := evt.Info.Sender
    
    // Check if message is from group
    isGroup := evt.Info.Chat.Server == "g.us"

    // Print message info
    lib.PrintIncomingMessage(evt)

    // Get message content
    msg := evt.Message.GetConversation()
    if msg == "" {
        // Check for caption in image/video/document
        if img := evt.Message.GetImageMessage(); img != nil {
            msg = img.GetCaption()
        } else if vid := evt.Message.GetVideoMessage(); vid != nil {
            msg = vid.GetCaption()
        } else if doc := evt.Message.GetDocumentMessage(); doc != nil {
            msg = doc.GetCaption()
        }
    }

    // Handle different message types
    switch {
    case evt.Message.GetStickerMessage() != nil:
        h.handleSticker(evt)
    case evt.Message.GetDocumentMessage() != nil:
        h.handleDocument(evt)
    case evt.Message.GetAudioMessage() != nil:
        h.handleAudio(evt)
    }

    // Process commands
    if strings.HasPrefix(msg, "!") || strings.HasPrefix(msg, "/") || strings.HasPrefix(msg, ".") {
        h.handleCommand(evt, msg[1:], sender, isGroup)
    }
}

// handleCommand processes bot commands
func (h *MessageHandler) handleCommand(evt *events.Message, cmd string, sender types.JID, isGroup bool) {
    // Split command and args
    args := strings.Fields(cmd)
    if len(args) == 0 {
        return
    }

    command := strings.ToLower(args[0])
    
    // Check cooldown
    if !h.checkCooldown(sender.String()) {
        lib.SendQuotedTextReply(h.client, evt, "⏳ Please wait before using another command")
        return
    }

    // Command aliases
    commandAliases := map[string]string{
        "p": "ping",
        "h": "help",
        "i": "info",
        // Add more aliases as needed
    }

    // Check for command alias
    if alias, exists := commandAliases[command]; exists {
        command = alias
    }

    // Group-specific command checks
    if isGroup {
        if !h.isCommandAllowedInGroup(command) {
            lib.SendQuotedTextReply(h.client, evt, "⚠️ This command cannot be used in groups")
            return
        }
    }

    // Log command usage
    log.Printf("[%s] Command %s used by %s in %v\n", 
        time.Now().Format("2006-01-02 15:04:05"),
        command, 
        sender.String(),
        evt.Info.Chat)

    // Check if command requires admin/owner privileges
    if isAdminCommand(command) && !isAdmin(sender) {
        lib.SendQuotedTextReply(h.client, evt, "⚠️ This command requires admin privileges")
        return
    }

    // Route command to appropriate handler in plugins package
    // Command handling is done in plugins/router.go
}

// Helper functions

func (h *MessageHandler) checkCooldown(user string) bool {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    const cooldownTime = 3 // seconds

    now := lib.GetCurrentTimestamp()
    lastUsed, exists := h.cooldowns[user]

    if !exists || now-lastUsed >= cooldownTime {
        h.cooldowns[user] = now
        return true
    }

    return false
}

func isAdminCommand(cmd string) bool {
    adminCommands := []string{
        "ban",
        "unban",
        "kick",
        "promote",
        "demote",
    }

    for _, c := range adminCommands {
        if cmd == c {
            return true
        }
    }
    return false
}

func isAdmin(sender types.JID) bool {
    // Add admin JIDs here
    admins := []string{
        "6285179855248@s.whatsapp.net", // Owner
    }

    senderStr := sender.String()
    for _, admin := range admins {
        if senderStr == admin {
            return true
        }
    }
    return false
}

func (h *MessageHandler) handleSticker(evt *events.Message) {
    // Handle sticker messages
    sticker := evt.Message.GetStickerMessage()
    if sticker != nil {
        // Process sticker
        log.Printf("Received sticker from %s\n", evt.Info.Sender)
    }
}

func (h *MessageHandler) handleDocument(evt *events.Message) {
    // Handle document messages
    doc := evt.Message.GetDocumentMessage()
    if doc != nil {
        // Process document
        log.Printf("Received document: %s from %s\n", doc.GetFileName(), evt.Info.Sender)
    }
}

func (h *MessageHandler) handleAudio(evt *events.Message) {
    // Handle audio messages
    audio := evt.Message.GetAudioMessage()
    if audio != nil {
        // Process audio
        log.Printf("Received audio from %s\n", evt.Info.Sender)
    }
}

func (h *MessageHandler) isCommandAllowedInGroup(cmd string) bool {
    // Add commands that should be disabled in groups
    disabledInGroups := []string{
        "broadcast",
        "spam",
        // Add more as needed
    }

    for _, disabled := range disabledInGroups {
        if cmd == disabled {
            return false
        }
    }
    return true
}

// Add method to lib/function.go:
func GetCurrentTimestamp() int64 {
    return time.Now().Unix()
}