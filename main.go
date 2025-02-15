package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    waLog "go.mau.fi/whatsmeow/util/log"

    "go.mau.fi/whatsmeow/proto/waE2E"
    "google.golang.org/protobuf/proto"

    "github.com/FahriAdison/Alya-Go/plugins"
    "github.com/FahriAdison/Alya-Go/lib"
)

var client *whatsmeow.Client

func main() {
    // Database setup
    dbLog := waLog.Stdout("Database", "ERROR", true)
    container, err := sqlstore.New("sqlite3", "file:whatsapp-session.db?_foreign_keys=off", dbLog)
    if err != nil {
	log.Fatalf("❌ Error creating store: %v", err)
    }

    deviceStore, err := container.GetFirstDevice()
    if err != nil {
	log.Fatalf("❌ Error getting device: %v", err)
    }

    client = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))

    // ✅ Connect before generating pairing code
    if err := client.Connect(); err != nil {
	log.Fatalf("❌ Error connecting to WhatsApp: %v", err)
    }

    // ✅ Generate Pairing Code
    if client.Store.ID == nil {
	fmt.Println("🔄 Generating pairing code...")

	pairingCode, err := client.PairPhone("639687312284", true, whatsmeow.PairClientChrome, "Chrome (Windows)")
	if err != nil {
	    log.Fatalf("❌ Error generating pairing code: %v", err)
	}

	fmt.Printf("📌 Pairing Code: %s\n", pairingCode)
	fmt.Println("✅ Enter this pairing code in your WhatsApp app to connect!")

	// Give time for pairing before proceeding
	time.Sleep(10 * time.Second)
    }

    // ✅ Send Online Status
    sendOnlineIndicator()

    // ✅ Message Event Handling
    client.AddEventHandler(func(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
	    lib.PrintIncomingMessage(v)
	    plugins.Handle(client, v)
	}
    })

    fmt.Println("🤖 Bot is running (Press CTRL+C to exit)")

    // ✅ Graceful Shutdown Handling
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    client.Disconnect()
}

// ✅ Notify Admin That Bot is Online
func sendOnlineIndicator() {
    adminJID := types.NewJID("6285179855248", "s.whatsapp.net") // Replace with your admin JID
    msg := &waE2E.Message{
	Conversation: proto.String("🤖 Bot is now online!"),
    }
    _, err := client.SendMessage(context.Background(), adminJID, msg)
    if err != nil {
	fmt.Println("⚠️ Failed to send online indicator:", err)
    }
}
