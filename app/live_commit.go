package app

import (
	"net/http"
	"io"
)

// Handle websocket live commit
func LiveCommit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "TODO websocket")
}
