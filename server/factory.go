package server

import (
	"errors"
	"github.com/xxlixin1993/LCS/configure"
)

// Start server
func StartServer() error {
	protocolType := configure.DefaultString("server.support", "http")

	switch protocolType {
	case "http":
		return Run()
	default:
		return errors.New("unknow server.support in configure")
	}

	return nil
}
