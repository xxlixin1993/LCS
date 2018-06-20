package server

import (
	"errors"
	"github.com/xxlixin1993/LCS/configure"
	"github.com/xxlixin1993/LCS/server/rtmp"
	"github.com/xxlixin1993/LCS/server/http"
)

// Start server
func StartServer() error {
	protocolType := configure.DefaultString("server.support", "http")

	switch protocolType {
	case "http":
		return http.Run()
	case "rtmp":
		return rtmp.StartRtmp()
	default:
		return errors.New("unknow server.support in configure")
	}

	return nil
}
