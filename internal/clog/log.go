package clog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelStrings = [...]string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

type Logger struct {
	mu         sync.Mutex
	level      LogLevel
	out        io.Writer
	timeFormat string
	showCaller bool
}

var std = New()

func New() *Logger {
	return &Logger{
		level:      INFO,
		out:        os.Stderr,
		timeFormat: "2006-01-02 15:04:05",
		showCaller: true,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

func (l *Logger) ShowCaller(show bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showCaller = show
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func (l *Logger) log(level LogLevel, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var builder strings.Builder

	if l.timeFormat != "" {
		builder.WriteString(time.Now().Format(l.timeFormat))
		builder.WriteString(" ")
	}

	builder.WriteString("[")
	builder.WriteString(levelStrings[level])
	builder.WriteString("] ")

	if l.showCaller {
		caller := getCaller()
		if caller != "" {
			builder.WriteString(caller)
			builder.WriteString(" - ")
		}
	}

	builder.WriteString(msg)
	builder.WriteString("\n")

	fmt.Fprint(l.out, builder.String())

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(DEBUG, fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.log(INFO, fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(WARN, fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.log(ERROR, fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(FATAL, fmt.Sprint(v...))
}

func SetLevel(level LogLevel)     { std.SetLevel(level) }
func SetOutput(w io.Writer)       { std.SetOutput(w) }
func SetTimeFormat(format string) { std.SetTimeFormat(format) }
func ShowCaller(show bool)        { std.ShowCaller(show) }
func Debug(v ...interface{})      { std.Debug(v...) }
func Info(v ...interface{})       { std.Info(v...) }
func Warn(v ...interface{})       { std.Warn(v...) }
func Error(v ...interface{})      { std.Error(v...) }
func Fatal(v ...interface{})      { std.Fatal(v...) }

func NewFileWriter(filename string) (io.Writer, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
