package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type LogType int

const (
	Debug LogType = iota
	Fatal
	Error
)

type Log struct {
	message string
	logType LogType
}

type Logger struct {
	inputChn chan Log
	file     *os.File
}

func NewLogger(debugMode bool, bufsize int) Logger {
	file, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now()), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open file for logger: %s\n", err)
	}

	if debugMode {
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		log.SetOutput(file)
	}

	logger := Logger{
		inputChn: make(chan Log, bufsize),
		file:     file,
	}

	go func() {
		log.Printf("logger-channel: opened\n")
		for msg := range logger.inputChn {
			switch msg.logType {
			case Error:
				log.Printf("error: %s", msg.message)
			case Fatal:
				log.Fatalf("fatal: %s", msg.message)
			case Debug:
				log.Printf("debug: %s", msg.message)
			}
		}
		log.Printf("logger-channel: closed\n")
	}()

	return logger
}

func (l *Logger) sendFatalF(format string, v ...any) {
	l.sendMessage(fmt.Sprintf(format, v...), Fatal)
}

func (l *Logger) sendErrorf(format string, v ...any) {
	l.sendMessage(fmt.Sprintf(format, v...), Error)
}

func (l *Logger) sendMessageF(format string, v ...any) {
	l.sendMessage(fmt.Sprintf(format, v...), Debug)
}

func (l *Logger) sendMessage(msg string, logType LogType) {
	l.inputChn <- Log{message: msg, logType: logType}
}

func (l *Logger) sendLog(log Log) {
	l.inputChn <- log
}

func (l *Logger) terminate() error {
	close(l.inputChn)
	return l.file.Close()
}
