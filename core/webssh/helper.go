package webssh

import (
	"anylinker/common/log"
	"anylinker/core/utils/resp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
	"time"
)



func handleError(c *fiber.Ctx, err error) bool {
	if err != nil {
		resp.JSON(c,resp.ErrInternalServer, err.Error())
		return true
	}
	return false
}

func WshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		log.Error("handler ws ERROR:",zap.Error(err))
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			log.Error("websocket writes control message failed:",zap.Error(err))

		}
		return true
	}
	return false
}
