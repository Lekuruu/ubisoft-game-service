package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

const ERROR int = 40
const WARNING int = 30
const INFO int = 20
const DEBUG int = 10
const VERBOSE int = 0

type Logger struct {
	logger *log.Logger
	name   string
	level  int
}

func CreateLogger(name string, level int) *Logger {
	l := log.New(os.Stdout, "", 0)
	return &Logger{
		logger: l,
		name:   name,
		level:  level,
	}
}

func (c *Logger) formatLogMessage(level, msg string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] - <%s> %s: %s", timestamp, c.name, level, msg)
}

func (c *Logger) Info(msg string) {
	if c.level > INFO {
		return
	}
	c.logger.Println(c.formatLogMessage("INFO", msg))
}

func (c *Logger) Error(msg string) {
	if c.level > ERROR {
		return
	}
	c.logger.Println(c.formatLogMessage("ERROR", msg))
}

func (c *Logger) Warning(msg string) {
	if c.level > WARNING {
		return
	}
	c.logger.Println(c.formatLogMessage("WARNING", msg))
}

func (c *Logger) Debug(msg string) {
	if c.level > DEBUG {
		return
	}
	c.logger.Println(c.formatLogMessage("DEBUG", msg))
}

func (c *Logger) Verbose(msg string) {
	if c.level > VERBOSE {
		return
	}
	c.logger.Println(c.formatLogMessage("VERBOSE", msg))
}
