package lib

import (
    "context"
    "fmt"
    "os"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "go.mau.fi/whatsmeow/types/events"
    "google.golang.org/protobuf/proto"
)

// SendQuotedTextReply sends a text reply quoting the user's original message.
func SendQuotedTextReply(client *whatsmeow.Client, evt *events.Message, replyText string) error {
    msg := &waE2E.Message{
	ExtendedTextMessage: &waE2E.ExtendedTextMessage{
	    Text: proto.String(replyText),
	    ContextInfo: &waE2E.ContextInfo{
		StanzaID:  proto.String(evt.Info.ID),
		Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
		QuotedMessage: evt.Message,
	    },
	},
    }

    // For groups, evt.Info.Chat is the group JID; if empty, fallback to Sender for direct messages.
    target := evt.Info.Chat
    if target.IsEmpty() {
	target = evt.Info.Sender
    }

    resp, err := client.SendMessage(context.Background(), target, msg)
    if err != nil {
	return fmt.Errorf("failed to send quoted text reply: %w", err)
    }
    fmt.Println("Quoted text reply sent at:", resp.Timestamp)
    return nil
}

// SendQuotedImageReply sends an image + caption as a quoted reply.
// imagePath is the local path to the image file, caption is the text to display under the image.
func SendQuotedImageReply(client *whatsmeow.Client, evt *events.Message, imagePath string, caption string) error {
    // Read the image from disk.
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
	return fmt.Errorf("failed to read image from %s: %w", imagePath, err)
    }

    // Upload the image to WhatsApp servers.
    uploadResp, err := client.Upload(context.Background(), imageData, whatsmeow.MediaImage)
    if err != nil {
	return fmt.Errorf("failed to upload image: %w", err)
    }

    // Build an ImageMessage referencing the user’s original message in ContextInfo.
    imageMsg := &waE2E.ImageMessage{
	Caption:  proto.String(caption),
	Mimetype:  proto.String("image/jpeg"),
	URL:      proto.String(uploadResp.URL),
	DirectPath: proto.String(uploadResp.DirectPath),
	MediaKey:     uploadResp.MediaKey,
	FileEncSHA256: uploadResp.FileEncSHA256,
	FileSHA256:  uploadResp.FileSHA256,
	FileLength:   proto.Uint64(uploadResp.FileLength),
	ContextInfo: &waE2E.ContextInfo{ // Add ContextInfo here, inside ImageMessage
	    StanzaID:  proto.String(evt.Info.ID),
	    Participant: proto.String(evt.Info.Sender.ToNonAD().String()),
	    QuotedMessage: evt.Message,
	},
    }
    msg := &waE2E.Message{
	ImageMessage: imageMsg,
    }


    target := evt.Info.Chat
    if target.IsEmpty() {
	target = evt.Info.Sender
    }

    resp, err := client.SendMessage(context.Background(), target, msg)
    if err != nil {
	return fmt.Errorf("failed to send quoted image reply: %w", err)
    }
    fmt.Println("Quoted image reply sent at:", resp.Timestamp)
    return nil
}
