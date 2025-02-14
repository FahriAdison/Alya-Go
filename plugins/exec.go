package plugins

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
    "reflect" // For basic reflection-based "eval" attempt
    "runtime"
    "strings"
    "time"

    "github.com/FahriAdison/Alya-Go/lib"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
)

// init registers exec commands
func init() {
    RegisterCommand("=>", ExecHandler) // Go code evaluation with output
    RegisterCommand(">", ExecHandler)  // Go code evaluation with output (alternative prefix)
    RegisterCommand("$", ExecHandler)  // Shell command with output
}

// Owner number (replace with your actual owner number JID)
const botOwnerNumber = "6285179855248@s.whatsapp.net" // Replace with your JID

// ExecHandler handles command execution and Go code evaluation.
// SECURITY WARNING: Go code evaluation is VERY DANGEROUS and should ONLY be used by the bot owner in a safe environment.
func ExecHandler(client *whatsmeow.Client, evt *events.Message) {
    text := evt.Message.GetConversation()
    if text == "" {
	return
    }

    sender := evt.Info.Sender.ToNonAD().String() // Get sender JID without AD suffix
    if sender != botOwnerNumber {
	lib.SendQuotedTextReply(client, evt, "⚠️ Access denied. Go code execution is only allowed for the bot owner.")
	return
    }


    prefix := ""
    commandText := ""

    if strings.HasPrefix(text, "=>") {
	prefix = "=>"
	commandText = strings.TrimSpace(text[2:])
    } else if strings.HasPrefix(text, ">") {
	prefix = ">"
	commandText = strings.TrimSpace(text[1:])
    } else if strings.HasPrefix(text, "$") {
	prefix = "$"
	commandText = strings.TrimSpace(text[1:])
    } else {
	return // Not an exec command
    }

    if commandText == "" {
	lib.SendQuotedTextReply(client, evt, "Error: Command/Code is empty.")
	return
    }

    var replyText string
    var err error

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Timeout for commands/code
    defer cancel()

    switch prefix {
    case "=>", ">": // Go Code Evaluation
	// VERY BASIC and UNSAFE attempt at Go code evaluation using reflection.
	// DO NOT USE THIS WITH UNTRUSTED INPUT.
	replyText = evaluateGoCode(commandText)


    case "$": // Shell command execution (existing functionality - kept as before)
	var cmd *exec.Cmd
	if prefix == "$" {
	    // Execute in shell (/bin/sh -c on Unix-like, cmd /C on Windows - should work in Termux)
	    shell := "/bin/sh" // Default shell for Unix-like systems (including Termux)
	    if runtime.GOOS == "windows" {
		shell = "cmd"
	    }
	    cmd = exec.CommandContext(ctx, shell, "-c", commandText)
	}


	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	output := stdout.String()
	errorOutput := stderr.String()

	if err != nil {
	    replyText = fmt.Sprintf("Error executing shell command:\n`\n%s\n`\nStderr:\n`\n%s\n`", err, errorOutput)
	} else {
	    if output != "" {
		replyText = fmt.Sprintf("Shell Command Output:\n`\n%s\n`", truncateOutput(output, 3800)) // Increased truncate limit to 3800
	    } else if errorOutput != "" {
		replyText = fmt.Sprintf("Shell command executed successfully, but output to stderr:\n`\n%s\n`", truncateOutput(errorOutput, 3800)) // Increased truncate limit to 3800
	    } else {
		replyText = "Shell command executed successfully, no output."
	    }
	}


    default:
	replyText = "Error: Invalid command prefix. Use =>, >, or $."
    }

    if replyText != "" {
	lib.SendQuotedTextReply(client, evt, replyText)
    } else if err != nil {
	lib.SendQuotedTextReply(client, evt, fmt.Sprintf("An unexpected error occurred: %v", err))
    }
}

// evaluateGoCode - VERY BASIC and UNSAFE Go code evaluation attempt.
// DO NOT USE WITH UNTRUSTED INPUT.  This is highly limited and for demonstration ONLY.
func evaluateGoCode(code string) string {
    defer func() {
	if r := recover(); r != nil {
	    fmt.Println("Panic during Go code evaluation:", r) // Log panic to console
	    // Print detailed panic info to console for debugging, but don't send to WhatsApp (potentially sensitive)
	    // debug.PrintStack() // Uncomment for stack trace if needed during testing

	}
    }()

    // Very limited example - only trying to evaluate basic expressions or function calls
    // THIS IS NOT A SECURE OR ROBUST "eval()" IMPLEMENTATION.
    code = strings.TrimSpace(code)
    if code == "" {
	return "Error: Go code snippet is empty."
    }

    // Example: Attempt to parse as a simple expression and use reflection to get its value.
    // This is VERY limited and will fail for most Go code.
    var resultValue reflect.Value
    var evalError error


    // **VERY LIMITED AND UNSAFE EVALUATION ATTEMPT.  FOR DEMONSTRATION ONLY.**
    // In a real secure scenario, you would NEVER do this directly.
    // This is just a placeholder to show the *idea* and is not robust.
    // You would need a much more sophisticated (and probably sandboxed) approach for actual code execution.
    // For simple examples, we are just trying to evaluate basic Go expressions or function calls.
    // Example:  "1 + 2", "len(\"hello\")", "strings.ToUpper(\"test\")" (assuming import "strings" is somehow available - which it's not in this simple example)

    // In this VERY LIMITED EXAMPLE, we are NOT actually dynamically compiling and executing Go code in a safe way.
    // This is just a placeholder that will likely fail for anything but the simplest expressions.
    // For a real "eval()" in Go, you would need to use "go/parser", "go/types", "go/build", "reflect" in a much more complex and secure manner.

    // Placeholder example - just try to "interpret" some very simple hardcoded examples.
    switch code {
    case "1 + 1":
	resultValue = reflect.ValueOf(1 + 1)
    case "2 * 3":
	resultValue = reflect.ValueOf(2 * 3)
    case "len(\"test\")":
	resultValue = reflect.ValueOf(len("test"))
    default:
	return "Error: Invalid or unsupported Go code snippet for evaluation.\n\nThis is a *very basic* and *unsafe* 'eval'-like attempt.\nOnly extremely simple expressions are (very partially) supported for demonstration."
    }


    if evalError != nil {
	return fmt.Sprintf("Error evaluating Go code:\n%v", evalError)
    }

    if resultValue.IsValid() {
	return fmt.Sprintf("Go Code Result:\n`\n%v\n`\n(Type: %s)", resultValue.Interface(), resultValue.Type())
    } else {
	return "Go Code evaluated, but no result or output."
    }


}


// truncateOutput limits the output string to a maximum length and adds a "..." suffix if truncated.
func truncateOutput(output string, maxLength int) string {
    if len(output) > maxLength {
	return output[:maxLength-3] + "..."
    }
    return output
}
