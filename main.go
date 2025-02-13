package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/mdp/qrterminal/v3"
    _ "github.com/mattn/go-sqlite3"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func main() {
    // Enable CGO requirements
    os.Setenv("CGO_ENABLED", "1")

    // Initialize logger
    dbLog := waLog.Stdout("Database", "ERROR", true)
    container, err := sqlstore.New("sqlite3", "file:whatsapp-session.db?_foreign_keys=off", dbLog)
    if err != nil {
	panic(err)
    }

    // Get device store
    deviceStore, err := container.GetFirstDevice()
    if err != nil {
	panic(err)
    }

    // Create client
    client = whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))

    // QR Code Login Flow
    if client.Store.ID == nil {
	// Generate QR channel
	qrChan, _ := client.GetQRChannel(context.Background())

	// Connect in background
	go func() {
	    if err := client.Connect(); err != nil {
		panic(err)
	    }
	}()

	// Display QR code
	fmt.Println("Waiting for QR code...")
	for evt := range qrChan {
	    if evt.Event == "code" {
		qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		fmt.Println("Scan this QR code with WhatsApp Mobile App")
	    } else if evt.Event == "success" {
		fmt.Println("Login successful!")
		break
	    }
	}
    } else {
	// Existing session
	if err := client.Connect(); err != nil {
	    panic(err)
	}
    }

    // Keep alive
    fmt.Println("Bot is running (Press CTRL+C to exit)")
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c
    client.Disconnect()
}