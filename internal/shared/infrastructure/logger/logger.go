package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/WhiteRain7/GoTestCase/internal/shared/infrastructure/config"
)

type msgCat string

const (
	msgInfo    msgCat = "INFO"
	msgWarning msgCat = "WARN"
	msgError   msgCat = "ERRO"
	msgFatal   msgCat = "FATL"
)

func (c msgCat) Validate() bool {
	return c == msgInfo || c == msgWarning || c == msgError || c == msgFatal
}

type msg struct {
	category msgCat
	text     string
}

type Logger struct {
	queue   chan msg
	file    *os.File
	ctx     context.Context
	cancel  context.CancelFunc
	console bool
}

func NewLogger(cfg config.LoggerConfig) (*Logger, error) {
	f, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	logger := &Logger{queue: make(chan msg), file: f, ctx: ctx, cancel: cancel, console: cfg.Console}
	logger.write(msg{msgInfo, "Logger start"})

	go func() {
		for entry := range logger.queue {
			logger.write(entry)
		}
		logger.file.Close()
	}()

	return logger, nil
}

func (l *Logger) write(s msg) {
	now := time.Now()
	msg := fmt.Sprintf("%s [%s] %s\n", now.Format("2006-01-02 15:04:05"), s.category, s.text)
	l.file.WriteString(msg)
	if l.console {
		fmt.Print(msg)
	}
}

func (l *Logger) log(s msg) {
	go func() {
		select {
		case l.queue <- s:
			return
		case <-l.ctx.Done():
			return
		}
	}()
}

func (l *Logger) Info(message string) {
	l.log(msg{msgInfo, message})
}

func (l *Logger) Warn(message string) {
	l.log(msg{msgWarning, message})
}

func (l *Logger) Error(message string) {
	l.log(msg{msgError, message})
}

func (l *Logger) Fatal(message string) {
	l.log(msg{msgFatal, message})
}

func (l *Logger) Close() {
	l.cancel()
	close(l.queue)
	l.write(msg{msgInfo, "Logger end"})
	l.file.Close()
}
