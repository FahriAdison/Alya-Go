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
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	// Local plugins
	_ "github.com/FahriAdison/Alya-Go/plugins"
)

var client *whatsmeow.Client

func main() {
	// Enable CGO
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
		qrChan, _ := client.GetQRChannel(context.Background())
		go func() {
			if err := client.Connect(); err != nil {
				panic(err)
			}
		}()

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
		if err := client.Connect(); err != nil {
			panic(err)
		}
	}

	// Send online indicator
	sendOnlineIndicator()

	// Load plugins
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			plugins.HandleMessage(client, v)
		}
	})

	// Keep alive
	fmt.Println("Bot is running (Press CTRL+C to exit)")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}

// Send online indicator to admin
func sendOnlineIndicator() {
	adminJID := "1234567890@s.whatsapp.net" // Replace with your number
	_, err := client.SendMessage(context.Background(), types.NewJID(adminJID), &whatsmeow.TextMessage{
		Content: "🤖 Bot is now online!",
	})
	if err != nil {
		fmt.Println("Failed to send online indicator:", err)
	}
}
