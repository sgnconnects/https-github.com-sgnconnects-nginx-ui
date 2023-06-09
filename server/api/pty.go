package api

import (
	"github.com/0xJacky/Nginx-UI/server/pkg/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func Pty(c *gin.Context) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// upgrade http to websocket
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("pty ws upgrade error", err)
		return
	}

	defer ws.Close()

	p, err := pty.NewPipeLine(ws)

	if err != nil {
		log.Println("pty.NewPipLine error", err)
		return
	}

	defer p.Pty.Close()

	errorChan := make(chan error, 1)
	go p.ReadPtyAndWriteWs(errorChan)
	go p.ReadWsAndWritePty(errorChan)

	err = <-errorChan

	if err != nil {
		log.Println(err)
	}

	return
}
