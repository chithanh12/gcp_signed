package main

import (
	"context"
	"github.com/chithanh12/gcp_signed/server"
	"os"
	"os/signal"
	"time"
)

func main() {
	s:= server.NewServer()
	s.Start()

	gracefulShutdown(s)
}

func gracefulShutdown(s *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
