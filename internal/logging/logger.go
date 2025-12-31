package logging

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Disable debug by default
	if os.Getenv("DEBUG") != "true" {
		DebugLogger.SetOutput(io.Discard)
	}
}

// RedactSensitive redacts sensitive information from strings
func RedactSensitive(s string) string {
	// Redact anything that looks like a token or secret
	if len(s) > 20 {
		return s[:8] + "..." + s[len(s)-4:]
	}
	return "***"
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}

// Fatal logs a fatal error and exits
func Fatal(format string, v ...interface{}) {
	ErrorLogger.Fatalf(format, v...)
}
