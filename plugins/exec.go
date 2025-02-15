package plugins

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
    "runtime"
    "strings"
    "time"

    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

const (
    maxExecTime     = 10 * time.Second // Match JS bot's 10s timeout
    maxOutputLength = 1500             // Limit output size
    ownerNumber     = "6285179855248"  // Replace with your number
)

// Register commands
func init() {
    RegisterCommand("$", ExecHandler)  // Shell commands
    RegisterCommand("=>", ExecHandler) // Go evaluation
    RegisterCommand(">", ExecHandler)  // Alternative eval prefix
}

// ExecHandler processes shell and code execution commands
func ExecHandler(client *whatsmeow.Client, evt *events.Message) {
    text := strings.TrimSpace(evt.Message.GetConversation())
    sender := evt.Info.Sender.String()

    // ✅ Restrict access to owner
    if !isOwner(sender) {
	lib.SendQuotedTextReply(client, evt, "⚠️ Akses ditolak! Hanya owner yang bisa menggunakan fitur ini!")
	return
    }

    // ✅ Detect command prefix
    prefix, cmd := parseCommand(text)

    var reply string
    var err error

    switch prefix {
    case "$":
	reply, err = executeShell(cmd)
    case "=>", ">":
	reply, err = evaluateCode(cmd)
    default:
	return
    }

    // ✅ Format error messages
    if err != nil {
	reply = fmt.Sprintf("❌ Error:\n```\n%s\n```", err.Error())
    }

    // ✅ Send response, truncating if necessary
    lib.SendQuotedTextReply(client, evt, truncateOutput(reply, maxOutputLength))
}

// ✅ Check if the sender is the bot owner
func isOwner(sender string) bool {
    ownerJID := ownerNumber + "@s.whatsapp.net"
    return sender == ownerJID
}

// ✅ Parse command prefix
func parseCommand(text string) (prefix, cmd string) {
    for _, p := range []string{"$ ", "=> ", "> "} {
	if strings.HasPrefix(text, p) {
	    return p[:len(p)-1], strings.TrimSpace(text[len(p):])
	}
    }
    return "", ""
}

// ✅ Execute shell commands
func executeShell(command string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), maxExecTime)
    defer cancel()

    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
	cmd = exec.CommandContext(ctx, "cmd", "/C", command)
    } else {
	cmd = exec.CommandContext(ctx, "sh", "-c", command)
    }

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    output := strings.TrimSpace(stdout.String())
    errorOutput := strings.TrimSpace(stderr.String())

    if err != nil {
	return "", fmt.Errorf("%s\n%s", errorOutput, err.Error())
    }

    if output == "" && errorOutput != "" {
	return fmt.Sprintf("⚠️ Warning:\n```\n%s\n```", errorOutput), nil
    }

    return fmt.Sprintf("💻 Output:\n```\n%s\n```", output), nil
}

// ✅ Evaluate Go expressions (Simple Implementation)
func evaluateCode(code string) (string, error) {
    switch {
    case strings.HasPrefix(code, "time."):
	return evalTimeExpression(code)
    case strings.HasPrefix(code, "strings."):
	return evalStringExpression(code)
    default:
	return "", fmt.Errorf("unsupported expression type")
    }
}

// ✅ Time-related evaluation
func evalTimeExpression(expr string) (string, error) {
    if expr == "time.Now()" {
	return time.Now().String(), nil
    }
    return "", fmt.Errorf("time expression not supported")
}

// ✅ String-related evaluation
func evalStringExpression(expr string) (string, error) {
    if expr == `strings.ToUpper("hello")` {
	return strings.ToUpper("hello"), nil
    }
    return "", fmt.Errorf("string operations not implemented")
}

// ✅ Truncate output to avoid message overflow
func truncateOutput(output string, max int) string {
    if len(output) > max {
	return output[:max-3] + "..."
    }
    return output
}
