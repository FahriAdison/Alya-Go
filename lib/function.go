package lib

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/gabriel-vasile/mimetype"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "google.golang.org/protobuf/proto"
)

// ########################################
// #         CORE MESSAGING FUNCTIONS     #
// ########################################

func GenerateMessageID(client *whatsmeow.Client) types.MessageID {
    return client.GenerateMessageID()
}

func SendText(client *whatsmeow.Client, target types.JID, text string, extra ...whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error) {
    return client.SendMessage(context.Background(), target, &waE2E.Message{
        Conversation: proto.String(text),
    }, extra...)
}

func SendQuotedTextReply(client *whatsmeow.Client, evt *events.Message, text string, extra ...whatsmeow.SendRequestExtra) error {
    target := evt.Info.Chat
    if target.IsEmpty() {
        target = evt.Info.Sender
    }

    _, err := client.SendMessage(context.Background(), target, &waE2E.Message{
        ExtendedTextMessage: &waE2E.ExtendedTextMessage{
            Text: proto.String(text),
            ContextInfo: &waE2E.ContextInfo{
                StanzaID:      proto.String(evt.Info.ID),
                Participant:   proto.String(evt.Info.Sender.ToNonAD().String()),
                QuotedMessage: evt.Message,
            },
        },
    }, extra...)
    return err
}

// ########################################
// #         MEDIA MESSAGING FUNCTIONS    #
// ########################################

func SendImage(client *whatsmeow.Client, evt *events.Message, imagePath string, caption string, extra ...whatsmeow.SendRequestExtra) error {
    return sendMediaMessage(client, evt, imagePath, caption, whatsmeow.MediaImage, extra...)
}

func SendDocument(client *whatsmeow.Client, evt *events.Message, filePath string, caption string, extra ...whatsmeow.SendRequestExtra) error {
    return sendMediaMessage(client, evt, filePath, caption, whatsmeow.MediaDocument, extra...)
}

func sendMediaMessage(client *whatsmeow.Client, evt *events.Message, filePath string, caption string, mediaType whatsmeow.MediaType, extra ...whatsmeow.SendRequestExtra) error {
    fileData, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    mime, err := mimetype.DetectFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to detect MIME type: %w", err)
    }

    uploadResp, err := client.Upload(context.Background(), fileData, mediaType)
    if err != nil {
        return fmt.Errorf("upload failed: %w", err)
    }

    target := evt.Info.Chat
    if target.IsEmpty() {
        target = evt.Info.Sender
    }

    var message *waE2E.Message
    switch mediaType {
    case whatsmeow.MediaImage:
        message = &waE2E.Message{
            ImageMessage: &waE2E.ImageMessage{
                Caption:       proto.String(caption),
                Mimetype:      proto.String(mime.String()),
                URL:           proto.String(uploadResp.URL),
                DirectPath:    proto.String(uploadResp.DirectPath),
                MediaKey:      uploadResp.MediaKey,
                FileEncSHA256: uploadResp.FileEncSHA256,
                FileSHA256:    uploadResp.FileSHA256,
                FileLength:    proto.Uint64(uploadResp.FileLength),
                ContextInfo: &waE2E.ContextInfo{
                    StanzaID:      proto.String(evt.Info.ID),
                    Participant:   proto.String(evt.Info.Sender.ToNonAD().String()),
                    QuotedMessage: evt.Message,
                },
            },
        }
    case whatsmeow.MediaDocument:
        message = &waE2E.Message{
            DocumentMessage: &waE2E.DocumentMessage{
                Caption:       proto.String(caption),
                Mimetype:      proto.String(mime.String()),
                FileName:      proto.String(filepath.Base(filePath)),
                URL:           proto.String(uploadResp.URL),
                DirectPath:    proto.String(uploadResp.DirectPath),
                MediaKey:      uploadResp.MediaKey,
                FileEncSHA256: uploadResp.FileEncSHA256,
                FileSHA256:    uploadResp.FileSHA256,
                FileLength:    proto.Uint64(uploadResp.FileLength),
                ContextInfo: &waE2E.ContextInfo{
                    StanzaID:      proto.String(evt.Info.ID),
                    Participant:   proto.String(evt.Info.Sender.ToNonAD().String()),
                    QuotedMessage: evt.Message,
                },
            },
        }
    }

    _, err = client.SendMessage(context.Background(), target, message, extra...)
    return err
}


// ########################################
// #         GROUP MANAGEMENT            #
// ########################################


func AddGroupParticipants(client *whatsmeow.Client, groupJID types.JID, users []types.JID) error {
    _, err := client.UpdateGroupParticipants(groupJID, users, whatsmeow.ParticipantChangeAdd)
    return err
}

// ########################################
// #         UTILITY FUNCTIONS           #
// ########################################


func SetDisappearingTimer(client *whatsmeow.Client, target types.JID, duration time.Duration) error {
    return client.SetDisappearingTimer(target, duration)
}