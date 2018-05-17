package server

import (
	"context"
	"github.com/xxlixin1993/LCS/app"
	"github.com/xxlixin1993/LCS/configure"
	"github.com/xxlixin1993/LCS/graceful_exit"
	"net/http"
	"time"
	"strconv"
	"io"
)

const KHttpServerModuleName = "httpServerModule"

var httpServer *HttpServer

type HttpServer struct {
	host       string
	port       string
	socketLink string
	server     *http.Server
}

// Implement ExitInterface
func (h *HttpServer) GetModuleName() string {
	return KHttpServerModuleName
}

// Implement ExitInterface
func (h *HttpServer) Stop() error {
	quitTimeout := configure.DefaultInt("http.quit_timeout", 30)
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(quitTimeout)*time.Second)

	return httpServer.server.Shutdown(ctx)
}

// Initialize http server
func initHttpServer() error {
	host := configure.DefaultString("host", "0.0.0.0")
	port := configure.DefaultString("port", "8080")
	readTimeout := configure.DefaultInt("http.read_timeout", 4)
	writeTimeout := configure.DefaultInt("http.write_timeout", 3)
	socketLink := host + ":" + port
	router := &RouterHandler{}

	httpServer = &HttpServer{
		host:       host,
		port:       port,
		socketLink: socketLink,
		server: &http.Server{
			Addr:         socketLink,
			Handler:      router,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
		},
	}

	// graceful exit
	if httpErr := graceful_exit.GetExitList().Pop(httpServer); httpErr != nil {
		return httpErr
	}

	return nil
}

// Start http server
func Run() error {
	initErr := initHttpServer()
	if initErr != nil {
		return initErr
	}

	serveErr := httpServer.server.ListenAndServe()
	if serveErr != nil {
		return serveErr
	}

	return nil
}

type RouterHandler struct {
}

// RouterHandler implements http.Handler.
func (rh *RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO Optimize later
	length := len(r.URL.Path)
	if length == 0 {
		io.WriteString(w, "Not supported, only roomId")
		return
	}

	roomIdInt, err := strconv.Atoi(r.URL.Path[1:length])
	if err != nil {
		io.WriteString(w, "Not supported, url path should be int")
		return
	}

	roomIdUint32 := uint32(roomIdInt)
	app.LiveCommit(w, r, roomIdUint32)
}
