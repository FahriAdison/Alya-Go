package plugins

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/FahriAdison/Alya-Go/lib" // ✅ Import lib package
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "google.golang.org/protobuf/proto"
)

// Register sticker command
func init() {
    RegisterCommand("sticker", StickerHandler)
    RegisterCommand("s", StickerHandler)
}

// StickerHandler converts images to stickers
func StickerHandler(client *whatsmeow.Client, evt *events.Message) {
    fmt.Println("[DEBUG] StickerHandler triggered")

    // Ensure message contains an image
    img := evt.Message.GetImageMessage()
    if img == nil {
	fmt.Println("[ERROR] No image found in the message")
	lib.SendQuotedTextReply(client, evt, "⚠️ Please send an *image* with the caption *sticker* or *s*.")
	return
    }

    // Check if caption is "sticker" or "s"
    caption := strings.ToLower(img.GetCaption())
    if caption != "sticker" && caption != "s" {
	fmt.Println("[ERROR] Caption does not match 'sticker' or 's'")
	return
    }

    fmt.Println("[DEBUG] Image detected with correct caption. Processing...")

    // Download image
    data, err := client.Download(img)
    if err != nil {
	fmt.Println("[ERROR] Failed to download image:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: Could not download the image.")
	return
    }

    // Create temp directory
    tempDir := filepath.Join("storage", "temp")
    if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
	fmt.Println("[ERROR] Failed to create temp directory:", err)
	return
    }

    // Save image to file
    inputPath := filepath.Join(tempDir, "input.jpg")
    err = os.WriteFile(inputPath, data, 0644)
    if err != nil {
	fmt.Println("[ERROR] Failed to save image:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: Could not save the image.")
	return
    }

    // Convert to WebP using ffmpeg
    outputPath := filepath.Join(tempDir, "output.webp")
    cmd := exec.Command("ffmpeg", "-i", inputPath, "-vf", "scale=512:512", "-lossless", "1", "-y", outputPath)
    fmt.Println("[DEBUG] Running ffmpeg command:", cmd.String())

    err = cmd.Run()
    if err != nil {
	fmt.Println("[ERROR] ffmpeg failed:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: *ffmpeg* failed to convert the image.")
	return
    }

    // Read WebP file
    webpData, err := os.ReadFile(outputPath)
    if err != nil {
	fmt.Println("[ERROR] Failed to read WebP file:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: Could not read the WebP file.")
	return
    }

    // Upload sticker
    uploadResp, err := client.Upload(context.Background(), webpData, whatsmeow.MediaImage)
    if err != nil {
	fmt.Println("[ERROR] Upload failed:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: Failed to upload sticker.")
	return
    }

    // Create StickerPackMessage
    stickerPack := &waE2E.StickerPackMessage{
	StickerPackID: proto.String("Alya-Go-Stickers"),
	Name:          proto.String("Alya-Go Stickers"),
	Publisher:     proto.String("Alya-Go Bot"),
	FileLength:    proto.Uint64(uint64(len(webpData))),
	FileSHA256:    uploadResp.FileSHA256,
	FileEncSHA256: uploadResp.FileEncSHA256,
	MediaKey:      uploadResp.MediaKey, // ✅ Fixed type issue
	DirectPath:    proto.String(uploadResp.DirectPath),
	ContextInfo: &waE2E.ContextInfo{
	    StanzaID:      proto.String(evt.Info.ID),
	    Participant:   proto.String(evt.Info.Sender.String()),
	    QuotedMessage: evt.Message,
	    RemoteJID:     proto.String(evt.Info.Chat.String()),
	},
    }

    // Send sticker
    _, err = client.SendMessage(context.Background(), evt.Info.Chat, &waE2E.Message{
	StickerMessage: &waE2E.StickerMessage{
	    URL:           proto.String(uploadResp.URL),
	    DirectPath:    proto.String(uploadResp.DirectPath),
	    MediaKey:      uploadResp.MediaKey, // ✅ Fixed type issue
	    Mimetype:      proto.String("image/webp"),
	    FileEncSHA256: uploadResp.FileEncSHA256,
	    FileSHA256:    uploadResp.FileSHA256,
	    FileLength:    proto.Uint64(uint64(len(webpData))),
	    ContextInfo:   stickerPack.ContextInfo,
	},
    })
    if err != nil {
	fmt.Println("[ERROR] Failed to send sticker:", err)
	lib.SendQuotedTextReply(client, evt, "❌ Error: Could not send the sticker.")
    }
}
