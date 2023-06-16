package emailClient

import (
	"log"
	"os"
)

var Logger *CustomLogger = NewLogger()

type CustomLogger struct {
	DebugLogger *log.Logger
	InfoLogger *log.Logger
	WarningLogger *log.Logger
}

func customLoggerComponent(level string) *log.Logger {
	return log.New(os.Stdout, " [" + level + "] ", log.Ldate | log.Ltime)
}

func NewLogger() *CustomLogger {
	return &CustomLogger{
		DebugLogger: customLoggerComponent("DEBUG"),
		InfoLogger: customLoggerComponent("INFO"),
		WarningLogger: customLoggerComponent("WARNING"),
	}
}

func (c *CustomLogger) Debug(msg any) {
	c.DebugLogger.Println(msg)
}

func (c *CustomLogger) Info(msg any) {
	c.InfoLogger.Println(msg)
}

func (c *CustomLogger) Warning(msg any) {
	c.WarningLogger.Println(msg)
}
