package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Printf("[INFO] %s %v", msg, args)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Printf("[ERROR] %s %v", msg, args)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Printf("[WARN] %s %v", msg, args)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		l.Printf("[DEBUG] %s %v", msg, args)
	}
}

func (l *Logger) WithContext(ctx string) *contextLogger {
	return &contextLogger{
		logger:  l,
		context: ctx,
	}
}

func (l *Logger) WithUser(userID int) *contextLogger {
	return &contextLogger{
		logger:  l,
		context: fmt.Sprintf("user=%d", userID),
	}
}

func (l *Logger) WithTenant(tenantID string) *contextLogger {
	return &contextLogger{
		logger:  l,
		context: fmt.Sprintf("tenant=%s", tenantID),
	}
}

type contextLogger struct {
	logger  *Logger
	context string
}

func (cl *contextLogger) Info(msg string, args ...interface{}) {
	cl.logger.Printf("[INFO] [%s] %s %v", cl.context, msg, args)
}

func (cl *contextLogger) Error(msg string, args ...interface{}) {
	cl.logger.Printf("[ERROR] [%s] %s %v", cl.context, msg, args)
}

func (cl *contextLogger) Warn(msg string, args ...interface{}) {
	cl.logger.Printf("[WARN] [%s] %s %v", cl.context, msg, args)
}

func (cl *contextLogger) Debug(msg string, args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		cl.logger.Printf("[DEBUG] [%s] %s %v", cl.context, msg, args)
	}
}
