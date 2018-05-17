package server

import (
	"github.com/xxlixin1993/LCS/app"
	"io"
	"net/http"
	"strconv"
)

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
