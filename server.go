package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/websocket"
)

// Upgrader ใช้ในการอัปเกรดการเชื่อมต่อ HTTP ให้เป็น WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ตัวจัดการ WebSocket
func handleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return err
	}
	defer ws.Close()

	for {
		// อ่านข้อความจาก client
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		log.Printf("Received: %s", msg)

		// ส่งข้อความกลับไปยัง client
		err = ws.WriteMessage(mt, msg)
		if err != nil {
			log.Println("Write Error:", err)
			break
		}
	}
	return nil
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ตั้งค่าเส้นทางสำหรับ WebSocket
	e.GET("/ws", handleWebSocket)

	// เริ่มต้นเซิร์ฟเวอร์
	e.Logger.Fatal(e.Start(":8080"))
}