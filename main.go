package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/mdp/qrterminal/v3"
    _ "github.com/mattn/go-sqlite3" // registers SQLite3 driver

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    waLog "go.mau.fi/whatsmeow/util/log"
    "go.mau.fi/whatsmeow/proto/waE2E"
    "google.golang.org/protobuf/proto"

    "github.com/FahriAdison/Alya-Go/plugins"
)

var client *whatsmeow.Client

func main() {
    // Initialize logger for the database.
    dbLog := waLog.Stdout("Database", "ERROR", true)

    // Create the store using the "sqlite3" driver.
    container, err := sqlstore.New("sqlite3", "file:whatsapp-session.db?_foreign_keys=off", dbLog)
    if err != nil {
	panic(fmt.Errorf("Error creating store: %w", err))
    }

    // Retrieve the first device.
    deviceStore, err := container.GetFirstDevice()
    if err != nil {
	panic(fmt.Errorf("Error getting device: %w", err))
    }

    // Create a new WhatsApp client.
    client = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))

    // QR Code Login Flow.
    if client.Store.ID == nil {
	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
	    panic(fmt.Errorf("Error getting QR channel: %w", err))
	}
	go func() {
	    if err := client.Connect(); err != nil {
		panic(err)
	    }
	}()

	fmt.Println("Waiting for QR code...")
	for evt := range qrChan {
	    if evt.Event == "code" {
		qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		fmt.Println("Scan this QR code with your WhatsApp mobile app")
	    } else if evt.Event == "success" {
		fmt.Println("Login successful!")
		break
	    }
	}
    } else {
	if err := client.Connect(); err != nil {
	    panic(fmt.Errorf("Error connecting: %w", err))
	}
    }

    // Send an online indicator.
    sendOnlineIndicator()

    // Use the command router: pass incoming messages to plugins.Handle.
    client.AddEventHandler(func(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
	    plugins.Handle(client, v)
	}
    })

    fmt.Println("Bot is running (Press CTRL+C to exit)")
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    client.Disconnect()
}

// sendOnlineIndicator sends a simple "Bot is now online" message to the admin.
func sendOnlineIndicator() {
    adminJID := types.NewJID("6285179855248", "s.whatsapp.net")
    msg := &waE2E.Message{
	Conversation: proto.String("🤖 Bot is now online!"),
    }
    _, err := client.SendMessage(context.Background(), adminJID, msg)
    if err != nil {
	fmt.Println("Failed to send online indicator:", err)
    }
}