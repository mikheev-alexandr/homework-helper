package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20, // 1Mb
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal: %v", sig)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdown()

	log.Println("Shutting down server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}

	log.Println("Server exeting")

	return nil
}
