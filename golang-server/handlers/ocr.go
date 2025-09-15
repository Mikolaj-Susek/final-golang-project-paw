package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/example/golang-postgres-crud/ocr"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "python-server:50051"
)

func PerformOcr(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	imageBytes := make([]byte, file.Size)
	_, err = fileContent.Read(imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Printf("Error connecting to gRPC server: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not connect to OCR server"})
		return
	}
	defer conn.Close()

	client := pb.NewOcrServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req := &pb.OcrRequest{ImageData: imageBytes}
	resp, err := client.PerformOcr(ctx, req)
	if err != nil {
		log.Printf("Error during OCR operation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not perform OCR operation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"extracted_text": resp.GetExtractedText()})
}
