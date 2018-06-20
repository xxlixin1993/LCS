package rtmp

import (
	"github.com/xxlixin1993/LCS/configure"
	"net"
	"github.com/xxlixin1993/LCS/logging"
	"github.com/xxlixin1993/LCS/server/rtmp/protocol"
	"log"
)

type Server struct {
}

func StartRtmp() error {
	host := configure.DefaultString("host", "0.0.0.0")
	port := configure.DefaultString("port", "1935")
	socketLink := host + ":" + port

	logging.TraceF("listen %s", socketLink)
	rtmpListen, err := net.Listen("tcp", socketLink)
	if err != nil {
		return err
	}

	defer rtmpListen.Close()
	rtmpServer := NewRtmpServer()

	rtmpServer.Serve(rtmpListen)
	return nil
}

func NewRtmpServer() *Server {
	return &Server{
	}
}

func (s *Server) Serve(listener net.Listener) error {
	for {
		netConn, err := listener.Accept()
		if err != nil {
			return err
		}

		conn := protocol.NewConn(netConn, 4*1024)
		logging.TraceF("new client, connect remote: %s, local:%s",
			netConn.RemoteAddr().String(), netConn.LocalAddr().String())

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *protocol.Conn) error {

	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Println("handleConn HandshakeServer err:", err)
		return err
	}

	// TODO 去读是否为connect命令
	return nil
}
