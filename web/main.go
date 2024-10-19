package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := &Server{logger}

	logger.Info("Starting server", "addr", *addr)
	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           server.routes(),
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
