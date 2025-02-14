package plugins

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// init registers menu commands
func init() {
    RegisterCommand("menu", MenuHandler)
    RegisterCommand("!menu", MenuHandler)
}

// MenuHandler sends the menu image with formatted caption
func MenuHandler(client *whatsmeow.Client, evt *events.Message) {
    caption := `📋 Bot Menu

Info

1. ping
2. menu

Owner

1. =>
2. >
3. $

Select an option by typing its command.`

    // Get current working directory
    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current working directory:", err)
        return
    }

    // Construct image path
    imagePath := filepath.Join(currentDir, "storage", "menu.jpg")

    // Verify image exists
    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        fmt.Printf("Error: Image file not found at %s\n", imagePath)
        fmt.Println("Please ensure 'storage/menu.jpg' exists in the project directory.")
        return
    }

    // Send image with caption
    err = lib.SendImage(client, evt, imagePath, caption)
    if err != nil {
        fmt.Println("Error sending menu image reply:", err)
    }
}