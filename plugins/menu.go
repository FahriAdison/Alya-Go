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

// MenuHandler calls the library function to send an image reply with a caption.
func MenuHandler(client *whatsmeow.Client, evt *events.Message) {
    caption := "📋 Bot Menu\n\n" +
	"1. /ping - Check connectivity\n" +
	"2. /menu - Show this menu\n\n" +
	"Select an option by typing its command."

    // Get the current working directory where the command is executed
    currentDir, err := os.Getwd()
    if err != nil {
	fmt.Println("Error getting current working directory:", err)
	return
    }

    // Construct the image path relative to the current working directory
    imagePath := filepath.Join(currentDir, "storage", "menu.jpg")

    // Check if the image file exists (for debugging)
    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
	fmt.Printf("Error: Image file not found at path: %s\n", imagePath)
	fmt.Println("Please ensure 'storage/menu.jpg' exists in the correct location relative to where you run 'go run main.go'.")
	return
    }

    err = lib.SendImage(client, evt, imagePath, caption)
    if err != nil {
	fmt.Println("Error sending menu image reply:", err)
    }
}
