package plugins

import (
    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    "strings"
    "fmt"
    "encoding/json"
    "net/http"
    "github.com/PuerkitoBio/goquery"
)

func init() {
    RegisterCommand("ai", AIHandler)
    RegisterCommand("dalle", ImageGenerationHandler)
    RegisterCommand("translate", TranslateHandler)
}

func AIHandler(client *whatsmeow.Client, evt *events.Message) {
    prompt := lib.GetMessageText(evt)
    if prompt == "" {
        lib.Reply(client, evt, "Please provide a prompt for the AI")
        return
    }
    
    // Using free GPT alternative API
    response, err := http.Get(fmt.Sprintf(
        "https://api.simsimi.net/v2/?text=%s&lc=id",
        url.QueryEscape(prompt),
    ))
    
    if err != nil {
        lib.Reply(client, evt, "Failed to get AI response")
        return
    }
    defer response.Body.Close()
    
    var result struct {
        Success string `json:"success"`
    }
    
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        lib.Reply(client, evt, "Failed to parse AI response")
        return
    }
    
    lib.Reply(client, evt, result.Success)
}

func ImageGenerationHandler(client *whatsmeow.Client, evt *events.Message) {
    prompt := lib.GetMessageText(evt)
    if prompt == "" {
        lib.Reply(client, evt, "Please provide a description for image generation")
        return
    }
    
    lib.Reply(client, evt, "Generating image... Please wait")
    
    // Using Unsplash API for relevant image search
    response, err := http.Get(fmt.Sprintf(
        "https://source.unsplash.com/1600x900/?%s",
        url.QueryEscape(prompt),
    ))
    
    if err != nil {
        lib.Reply(client, evt, "Failed to generate image")
        return
    }
    defer response.Body.Close()
    
    imageData, err := io.ReadAll(response.Body)
    if err != nil {
        lib.Reply(client, evt, "Failed to read image data")
        return
    }
    
    lib.SendImage(client, evt, imageData, prompt)
}

func TranslateHandler(client *whatsmeow.Client, evt *events.Message) {
    args := lib.ParseArgs(lib.GetMessageText(evt))
    if len(args) < 2 {
        lib.Reply(client, evt, "Usage: .translate <targetLang> <text>")
        return
    }
    
    targetLang := args[0]
    text := strings.Join(args[1:], " ")
    
    // Using LibreTranslate API
    payload := map[string]interface{}{
        "q": text,
        "source": "auto",
        "target": targetLang,
    }
    
    jsonData, _ := json.Marshal(payload)
    
    response, err := http.Post(
        "https://libretranslate.de/translate",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    
    if err != nil {
        lib.Reply(client, evt, "Failed to translate text")
        return
    }
    defer response.Body.Close()
    
    var result struct {
        TranslatedText string `json:"translatedText"`
    }
    
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        lib.Reply(client, evt, "Failed to parse translation")
        return
    }
    
    lib.Reply(client, evt, fmt.Sprintf("Translation (%s):\n%s", targetLang, result.TranslatedText))
}
