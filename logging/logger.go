package logging

import (
	"errors"
	"fmt"
	"path"
	"runtime"

	"github.com/xxlixin1993/CacheGo/utils"
	"github.com/xxlixin1993/LCS/configure"
	"github.com/xxlixin1993/LCS/graceful_exit"
	"sync"
)

// Log message level
const (
	KLevelFatal = iota
	KLevelError
	KLevelWarnning
	KLevelNotice
	KLevelInfo
	KLevelTrace
	KLevelDebug
)

// Log module name
const KLogModuleName = "logModule"

// Log output message level abbreviation
var LevelName = [7]string{"F", "E", "W", "N", "I", "T", "D"}

// Log instance
var loggerInstance *LogBase

// Log output type
const (
	KOutputFile   = "file"
	KOutputStdout = "stdout"
)

// Log interface. Need to be implemented when you want to extend.
type ILog interface {
	// Initialize Logger
	Init(config interface{}) error

	// Output message to log
	OutputLogMsg(msg []byte) error
}

// Log core program
type LogBase struct {
	mu sync.Mutex
	sync.WaitGroup
	handle  ILog
	message chan []byte
	skip    int
	level   int
}

// Initialize Log
func InitLog() error {
	outputType := configure.DefaultString("log.type", KOutputStdout)
	level := configure.DefaultInt("log.level", KLevelDebug)

	logger, err := createLogger(outputType, level)
	if err != nil {
		return err
	}

	// graceful exit
	graceful_exit.GetExitList().Pop(logger)

	go logger.Run()

	return err
}

// Implement ExitInterface
func (l *LogBase) GetModuleName() string {
	return KLogModuleName
}

// Implement ExitInterface
func (l *LogBase) Stop() error {
	close(loggerInstance.message)
	loggerInstance.Wait()
	return nil
}

// Create Logger instance
func createLogger(outputType string, level int) (*LogBase, error) {
	switch outputType {
	case KOutputStdout:
		loggerInstance = &LogBase{
			handle:  NewStdoutLog(),
			message: make(chan []byte, 1000),
			skip:    3,
			level:   level,
		}
		return loggerInstance, nil
	case KOutputFile:
		// TODO
		return nil, errors.New("TODO not supported")
	default:
		return nil, errors.New(configure.KUnknownTypeMsg)
	}
}

// Get Logger instance
func GetLogger() *LogBase {
	return loggerInstance
}

// Receive information, wait information
func (l *LogBase) Run() {
	loggerInstance.Add(1)

	for {
		msg, ok := <-l.message
		if !ok {
			l.Done()
			break
		}
		err := l.handle.OutputLogMsg(msg)
		if err != nil {
			fmt.Printf("Log: Output handle fail, err:%v\n", err.Error())
		}
	}
}

// Output message
func (l *LogBase) Output(nowLevel int, msg string) {
	now := utils.GetMicTimeFormat()

	l.mu.Lock()
	defer l.mu.Unlock()

	if nowLevel <= l.level {
		_, file, line, ok := runtime.Caller(l.skip)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = fmt.Sprintf("[%s] [%s %s:%d] %s\n", LevelName[nowLevel], now, filename, line, msg)
	}

	l.message <- []byte(msg)
}

func Debug(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelDebug, msg)
}

func DebugF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelDebug, msg)
}

func Trace(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelTrace, msg)
}

func TraceF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelTrace, msg)
}

func Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelInfo, msg)
}

func InfoF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelInfo, msg)
}
func Notice(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelNotice, msg)
}

func NoticeF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelNotice, msg)
}

func Warning(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelWarnning, msg)
}

func WarningF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelWarnning, msg)
}

func Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelError, msg)
}

func ErrorF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelError, msg)
}

func Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	GetLogger().Output(KLevelFatal, msg)
}

func FatalF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	GetLogger().Output(KLevelFatal, msg)
}
