package logging

import "fmt"

type LogStdout struct {
}

func NewStdoutLog() ILog {
	stdout := &LogStdout{}
	return stdout
}

// Init stdout
func (s *LogStdout) Init(config interface{}) error {
	return nil
}

// Output message in stdout.
func (s *LogStdout) OutputLogMsg(msg []byte) error {
	fmt.Print(string(msg))
	return nil
}
