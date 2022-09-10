package handlers

import (
	"net/http"
	"tinderutf/api/websocket"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

func Subscribe(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader := ws.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		id := c.GetString("userId")
		if id == "" {
			c.AbortWithStatusJSON(500, createDefaultError("invalid userId"))
			return
		}


		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		cl := websocket.NewClient(conn, hub)
		hub.InsertClient(cl, id)
	}
}