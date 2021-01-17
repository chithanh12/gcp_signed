package server

import (
	"context"
	"github.com/chithanh12/gcp_signed/signer"
	"github.com/labstack/echo/v4"
)

const (
	googleConfigFile ="poke-map-32809-a5632aa70eac.json"
	uploadBucket     = "signed-bucket-sample"
)

type Server struct {
	signer *signer.GcpSigner
	e *echo.Echo
}

func NewServer() *Server {
	s := &Server{
		e: echo.New(),
		signer: signer.NewGcpSigner(googleConfigFile, uploadBucket),
	}

	s.e.POST("/signed-upload",s.SignedUpload)
	s.e.POST("/signed-get", s.SignedGet)

	return s
}

func (s *Server) Start() {
	go func() {
		s.e.Start(":8080")
	}()
}

func (s *Server) Shutdown (ctx context.Context){
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Fatal(err)
	}
}


