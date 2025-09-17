package ocr

import (
	"context"
	"log"
	"time"

	pb "github.com/example/golang-postgres-crud/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "python-server:50051"
)

type OcrService struct {
	client pb.OcrServiceClient
	conn   *grpc.ClientConn
}

func NewOcrService() (*OcrService, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := pb.NewOcrServiceClient(conn)
	return &OcrService{
		client: client,
		conn:   conn,
	}, nil
}

func (s *OcrService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *OcrService) PerformOcr(imageBytes []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	req := &pb.OcrRequest{ImageData: imageBytes}
	resp, err := s.client.PerformOcr(ctx, req)
	if err != nil {
		log.Printf("Error during OCR operation: %v", err)
		return "", err
	}

	return resp.GetExtractedText(), nil
}
