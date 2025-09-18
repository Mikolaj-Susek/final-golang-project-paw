package handlers

import (
	"log"
	"net/http"

	"github.com/example/golang-postgres-crud/auth"
	ocr "github.com/example/golang-postgres-crud/ocr_service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TextReadingWebSocketHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not provided"})
		return
	}

	if err := auth.VerifyToken(tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token", "details": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}

		if messageType != websocket.BinaryMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("Only binary data (images) is supported."))
			continue
		}

		ocrService, err := ocr.NewOcrService()
		if err != nil {
			log.Printf("Error connecting to gRPC server: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Could not connect to OCR server"))
			return
		}
		defer ocrService.Close()

		ocrText, err := ocrService.PerformOcr(p)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Could not perform OCR operation"))
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(ocrText))
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}
