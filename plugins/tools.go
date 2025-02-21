package plugins

import (
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
