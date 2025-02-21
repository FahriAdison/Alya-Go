package plugins

import (
    "fmt"
    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

func init() {
    RegisterCommand("kick", KickHandler)
    RegisterCommand("add", AddHandler)
    RegisterCommand("promote", PromoteHandler)
    RegisterCommand("demote", DemoteHandler)
    RegisterCommand("setdesc", SetDescHandler)
}

// Add group management handlers

func KickHandler(client *whatsmeow.Client, evt *events.Message) {
    if !lib.IsGroup(evt) {
        lib.Reply(client, evt, "This command can only be used in groups")
        return
    }
    
    if !lib.IsAdmin(client, evt) {
        lib.Reply(client, evt, "You need to be an admin to use this command")
        return
    }
    
    mentioned := lib.GetMentionedUsers(evt)
    if len(mentioned) == 0 {
        lib.Reply(client, evt, "Please mention the user to kick")
        return
    }
    
    err := lib.KickGroupParticipant(client, evt.Chat.JID, mentioned[0])
    if err != nil {
        lib.Reply(client, evt, "Failed to kick user")
        return
    }
    
    lib.Reply(client, evt, "User has been kicked from the group")
}

func AddHandler(client *whatsmeow.Client, evt *events.Message) {
    if !lib.IsGroup(evt) {
        lib.Reply(client, evt, "This command can only be used in groups")
        return
    }
    
    numbers := lib.ParsePhoneNumbers(lib.GetMessageText(evt))
    if len(numbers) == 0 {
        lib.Reply(client, evt, "Please provide phone numbers to add")
        return
    }
    
    added := 0
    for _, number := range numbers {
        err := lib.AddGroupParticipant(client, evt.Chat.JID, number)
        if err == nil {
            added++
        }
    }
    
    lib.Reply(client, evt, fmt.Sprintf("Successfully added %d participants", added))
}

func PromoteHandler(client *whatsmeow.Client, evt *events.Message) {
    if !lib.IsGroup(evt) || !lib.IsAdmin(client, evt) {
        lib.Reply(client, evt, "Insufficient permissions")
        return
    }
    
    mentioned := lib.GetMentionedUsers(evt)
    if len(mentioned) == 0 {
        lib.Reply(client, evt, "Please mention the user to promote")
        return
    }
    
    err := lib.PromoteGroupParticipant(client, evt.Chat.JID, mentioned[0])
    if err != nil {
        lib.Reply(client, evt, "Failed to promote user")
        return
    }
    
    lib.Reply(client, evt, "User has been promoted to admin")
}

func DemoteHandler(client *whatsmeow.Client, evt *events.Message) {
    if !lib.IsGroup(evt) || !lib.IsAdmin(client, evt) {
        lib.Reply(client, evt, "Insufficient permissions")
        return
    }
    
    mentioned := lib.GetMentionedUsers(evt)
    if len(mentioned) == 0 {
        lib.Reply(client, evt, "Please mention the user to demote")
        return
    }
    
    err := lib.DemoteGroupParticipant(client, evt.Chat.JID, mentioned[0])
    if err != nil {
        lib.Reply(client, evt, "Failed to demote user")
        return
    }
    
    lib.Reply(client, evt, "User has been demoted from admin")
}

func SetDescHandler(client *whatsmeow.Client, evt *events.Message) {
    if !lib.IsGroup(evt) || !lib.IsAdmin(client, evt) {
        lib.Reply(client, evt, "Insufficient permissions")
        return
    }
    
    desc := lib.GetMessageText(evt)
    if desc == "" {
        lib.Reply(client, evt, "Please provide a new description")
        return
    }
    
    err := lib.SetGroupDescription(client, evt.Chat.JID, desc)
    if err != nil {
        lib.Reply(client, evt, "Failed to update group description")
        return
    }
    
    lib.Reply(client, evt, "Group description updated successfully")
}
