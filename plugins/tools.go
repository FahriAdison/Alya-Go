package plugins

import (
    "fmt"
    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

func init() {
    RegisterCommand("qr", QRGeneratorHandler)
    RegisterCommand("ss", ScreenshotHandler)
    RegisterCommand("weather", WeatherHandler)
    RegisterCommand("short", URLShortenerHandler)
}

// Add utility tool handlers

func QRGeneratorHandler(client *whatsmeow.Client, evt *events.Message) {
    text := lib.GetMessageText(evt)
    if text == "" {
        lib.Reply(client, evt, "Please provide text to generate QR code")
        return
    }
    
    qrCode, err := lib.GenerateQR(text)
    if (err != nil) {
        lib.Reply(client, evt, "Failed to generate QR code")
        return
    }
    
    lib.SendImage(client, evt, qrCode, "QR Code")
}

func ScreenshotHandler(client *whatsmeow.Client, evt *events.Message) {
    url := lib.GetMessageText(evt)
    if url == "" {
        lib.Reply(client, evt, "Please provide a URL to screenshot")
        return
    }
    
    screenshot, err := lib.CaptureWebsite(url)
    if err != nil {
        lib.Reply(client, evt, "Failed to capture screenshot")
        return
    }
    
    lib.SendImage(client, evt, screenshot, "Screenshot")
}

func WeatherHandler(client *whatsmeow.Client, evt *events.Message) {
    location := lib.GetMessageText(evt)
    if location == "" {
        lib.Reply(client, evt, "Please provide a location")
        return
    }
    
    weather, err := lib.GetWeatherInfo(location)
    if err != nil {
        lib.Reply(client, evt, "Failed to get weather information")
        return
    }
    
    response := fmt.Sprintf("Weather in %s:\nTemperature: %.1f°C\nCondition: %s\nHumidity: %d%%",
        weather.Location, weather.Temperature, weather.Condition, weather.Humidity)
    lib.Reply(client, evt, response)
}

func URLShortenerHandler(client *whatsmeow.Client, evt *events.Message) {
    url := lib.GetMessageText(evt)
    if url == "" {
        lib.Reply(client, evt, "Please provide a URL to shorten")
        return
    }
    
    shortURL, err := lib.ShortenURL(url)
    if err != nil {
        lib.Reply(client, evt, "Failed to shorten URL")
        return
    }
    
    lib.Reply(client, evt, fmt.Sprintf("Shortened URL: %s", shortURL))
}
