package plugins

import (
    "fmt"
    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    "github.com/gocolly/colly/v2"
    "github.com/kkdai/youtube/v2"
    "mvdan.cc/xurls/v2"
    "regexp"
)

func init() {
    RegisterCommand("yt", YoutubeHandler)
    RegisterCommand("tiktok", TiktokHandler) 
    RegisterCommand("ig", InstagramHandler)
}

func YoutubeHandler(client *whatsmeow.Client, evt *events.Message) {
    url := lib.GetMessageText(evt)
    if !lib.IsYoutubeURL(url) {
        lib.Reply(client, evt, "Please provide a valid YouTube URL")
        return
    }
    
    lib.Reply(client, evt, "Processing YouTube video...")
    
    // Using kkdai/youtube library for free YouTube downloads
    yt := youtube.Client{}
    video, err := yt.GetVideo(url)
    if err != nil {
        lib.Reply(client, evt, "Failed to get video information")
        return
    }
    
    // Get highest quality format under 50MB
    format := lib.GetOptimalFormat(video.Formats)
    if format == nil {
        lib.Reply(client, evt, "No suitable video format found")
        return
    }
    
    videoData, err := yt.DownloadFormat(video, format)
    if err != nil {
        lib.Reply(client, evt, "Failed to download video")
        return
    }
    
    caption := fmt.Sprintf("Title: %s\nDuration: %s\nChannel: %s",
        video.Title, video.Duration, video.Author)
    
    lib.SendVideo(client, evt, videoData, caption)
}

func TiktokHandler(client *whatsmeow.Client, evt *events.Message) {
    url := lib.GetMessageText(evt)
    if !lib.IsTiktokURL(url) {
        lib.Reply(client, evt, "Please provide a valid TikTok URL")
        return
    }
    
    lib.Reply(client, evt, "Processing TikTok video...")
    
    // Using web scraping with colly
    c := colly.NewCollector()
    var videoURL string
    
    c.OnHTML("video[src]", func(e *colly.HTMLElement) {
        videoURL = e.Attr("src")
    })
    
    err := c.Visit(url)
    if err != nil {
        lib.Reply(client, evt, "Failed to fetch TikTok video")
        return
    }
    
    if videoURL == "" {
        // Fallback to TikTok API
        videoURL, err = lib.GetTikTokNoWM(url)
        if err != nil {
            lib.Reply(client, evt, "Failed to get video URL")
            return
        }
    }
    
    videoData, err := lib.DownloadFile(videoURL)
    if err != nil {
        lib.Reply(client, evt, "Failed to download video")
        return
    }
    
    lib.SendVideo(client, evt, videoData, "TikTok Video")
}

func InstagramHandler(client *whatsmeow.Client, evt *events.Message) {
    url := lib.GetMessageText(evt)
    if !lib.IsInstagramURL(url) {
        lib.Reply(client, evt, "Please provide a valid Instagram URL")
        return
    }
    
    lib.Reply(client, evt, "Processing Instagram content...")
    
    // Using web scraping approach
    c := colly.NewCollector()
    var mediaURLs []string
    
    c.OnHTML(`meta[property="og:video"]`, func(e *colly.HTMLElement) {
        mediaURLs = append(mediaURLs, e.Attr("content"))
    })
    
    c.OnHTML(`meta[property="og:image"]`, func(e *colly.HTMLElement) {
        mediaURLs = append(mediaURLs, e.Attr("content"))
    })
    
    err := c.Visit(url)
    if err != nil {
        lib.Reply(client, evt, "Failed to fetch Instagram content")
        return
    }
    
    for _, mediaURL := range mediaURLs {
        mediaData, err := lib.DownloadFile(mediaURL)
        if err != nil {
            continue
        }
        
        if strings.Contains(mediaURL, ".mp4") {
            lib.SendVideo(client, evt, mediaData, "Instagram Video")
        } else {
            lib.SendImage(client, evt, mediaData, "Instagram Image")
        }
    }
}
