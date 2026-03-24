package handlers

import (
	"log"
	"looky/internal/utils"
	"looky/internal/ws"
	"strings"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

var upgrader = websocket.FastHTTPUpgrader{
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

func WebSocketHandler(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		auth := c.Get("Authorization")
		token = strings.TrimPrefix(auth, "Bearer ")
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	customerID := claims.UserID

	err = upgrader.Upgrade(c.RequestCtx(), func(conn *websocket.Conn) {
		client := ws.NewClient(conn)
		ws.H.Register(customerID, client)
		go client.WritePump()
		client.ReadPump(customerID)
	})

	if err != nil {
		log.Println("upgrade error:", err)
	}
	return err
}
