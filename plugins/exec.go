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
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
)

const (
    maxExecTime     = 15 * time.Second
    maxOutputLength = 1500
    ownerNumber     = "6285179855248" // Replace with your number
)

func init() {
    RegisterCommand("$", ExecHandler)  // Shell commands
    RegisterCommand("=>", ExecHandler) // Go evaluation
    RegisterCommand(">", ExecHandler)  // Alternative eval prefix
}

func ExecHandler(client *whatsmeow.Client, evt *events.Message) {
    if !isOwner(evt.Info.Sender) {
        lib.SendQuotedTextReply(client, evt, "⚠️ Akses ditolak! Hanya owner yang bisa menggunakan fitur ini!")
        return
    }

    text := strings.TrimSpace(evt.Message.GetConversation())
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

    if err != nil {
        reply = fmt.Sprintf("❌ Error:\n```\n%s\n```", err.Error())
    }
    
    lib.SendQuotedTextReply(client, evt, truncateOutput(reply, maxOutputLength))
}

func isOwner(sender types.JID) bool {
    ownerJID, _ := types.ParseJID(ownerNumber + "@s.whatsapp.net")
    return sender.ToNonAD() == ownerJID
}

func parseCommand(text string) (prefix, cmd string) {
    for _, p := range []string{"$ ", "=> ", "> "} {
        if strings.HasPrefix(text, p) {
            return p[:len(p)-1], strings.TrimSpace(text[len(p):])
        }
    }
    return "", ""
}

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

func evaluateCode(code string) (string, error) {
    // Remove this unused block:
    // allowed := map[string]interface{}{
    //     "math":   struct{}{},
    //     "time":   struct{}{},
    //     "strings": struct{}{},
    // }

    // Basic expression evaluation
    switch {
    case strings.HasPrefix(code, "math."):
        return evalMathExpression(code)
    case strings.HasPrefix(code, "time."):
        return evalTimeExpression(code)
    case strings.HasPrefix(code, "strings."):
        return evalStringExpression(code)
    default:
        return "", fmt.Errorf("unsupported expression type")
    }
}

func evalMathExpression(expr string) (string, error) {
    // Implement safe math evaluations
    return "", fmt.Errorf("math evaluation not implemented")
}

func evalTimeExpression(expr string) (string, error) {
    // Example: time.Now().Format("2006-01-02")
    if expr == "time.Now()" {
        return time.Now().String(), nil
    }
    return "", fmt.Errorf("time expression not supported")
}

func evalStringExpression(expr string) (string, error) {
    // Example: strings.ToUpper("hello")
    return "", fmt.Errorf("string operations not implemented")
}

func truncateOutput(output string, max int) string {
    if len(output) > max {
        return output[:max-3] + "..."
    }
    return output
}