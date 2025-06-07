package config

import (
	"log"
	"os"
	"time"
)

func InitLogger() {
	log.SetFlags(0) // Remove default flags
	log.SetOutput(os.Stdout)
	log.SetPrefix("") // Remove default prefix

	// Create a custom logger that formats time as [YYYY-MM-DD HH:MM:SS]
	log.SetOutput(&logWriter{
		writer: os.Stdout,
		format: "2006-01-02 15:04:05",
	})
}

// logWriter is a custom writer that formats the log output
type logWriter struct {
	writer *os.File
	format string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	// Add timestamp in the desired format
	timestamp := time.Now().Format(w.format)
	formatted := "[" + timestamp + "] " + string(p)
	return w.writer.Write([]byte(formatted))
}