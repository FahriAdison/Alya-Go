package lib

import "log"

// ErrorSeverity defines the severity level of an error
type ErrorSeverity int

const (
    SeverityLow ErrorSeverity = iota
    SeverityMedium
    SeverityHigh
    SeverityCritical
)

// HandleError processes errors based on severity
func HandleError(err error, severity ErrorSeverity, context string) {
    if err == nil {
        return
    }

    logMessage := fmt.Sprintf("[%s] Error: %v", context, err)

    switch severity {
    case SeverityCritical:
        log.Fatalf("CRITICAL: %s", logMessage)
    case SeverityHigh:
        log.Printf("HIGH: %s", logMessage)
    case SeverityMedium:
        log.Printf("MEDIUM: %s", logMessage)
    case SeverityLow:
        log.Printf("LOW: %s", logMessage)
    }
}
