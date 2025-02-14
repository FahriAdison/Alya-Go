package plugins

import (
    "fmt"
    "runtime"
    "strings"

    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// init registers ping commands
func init() {
    RegisterCommand("ping", PingHandler)
    RegisterCommand("!ping", PingHandler)
}

// PingHandler uses the library function to send a quoted system specifications reply.
func PingHandler(client *whatsmeow.Client, evt *events.Message) {
    systemSpecs := getSystemSpecs()
    err := lib.SendQuotedTextReply(client, evt, systemSpecs)
    if err != nil {
	fmt.Println("Error sending quoted system specs reply:", err)
    }
}

// getSystemSpecs gathers and formats system specification information.
func getSystemSpecs() string {
    var sb strings.Builder
    sb.WriteString("⚙️ *System Specifications*\n\n")

    // Operating System
    sb.WriteString("*OS:* ")
    sb.WriteString(runtime.GOOS)
    sb.WriteString("\n")

    // Architecture
    sb.WriteString("*Architecture:* ")
    sb.WriteString(runtime.GOARCH)
    sb.WriteString("\n")

    // Go Version
    sb.WriteString("*Go Version:* ")
    sb.WriteString(runtime.Version())
    sb.WriteString("\n")

    // CPU Cores
    sb.WriteString("*CPU Cores:* ")
    sb.WriteString(fmt.Sprintf("%d", runtime.NumCPU()))
    sb.WriteString(" Cores\n")

    // Memory Information (using runtime.MemStats)
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    sb.WriteString("*RAM Usage (Approximate):*\n")
    sb.WriteString("  - Alloc: ")
    sb.WriteString(formatBytes(m.Alloc)) // Bytes allocated and still in use
    sb.WriteString("\n")
    sb.WriteString("  - TotalAlloc: ")
    sb.WriteString(formatBytes(m.TotalAlloc)) // Total bytes allocated (even if freed)
    sb.WriteString("\n")
    sb.WriteString("  - Sys: ")
    sb.WriteString(formatBytes(m.Sys)) // Bytes of memory obtained from the OS
    sb.WriteString("\n")
    sb.WriteString("  - GcPause: ")
    sb.WriteString(fmt.Sprintf("%.4f ms", float64(m.PauseNs[(m.NumGC+299)%256])/1000000)) // Most recent pause in garbage collection
    sb.WriteString("\n")


    sb.WriteString("\n_Note: RAM usage is approximate and Go-specific._")

    return sb.String()
}

// formatBytes converts bytes to human-readable format (KB, MB, GB).
func formatBytes(b uint64) string {
    const unit = 1024
    if b < unit {
	return fmt.Sprintf("%d B", b)
    }
    div, exp := uint64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
	div *= unit
	exp++
    }
    return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
