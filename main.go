package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"pastebox.mohika.ir/internal/database"
	"pastebox.mohika.ir/web"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	dsn := flag.String("dsn", ":memory:", "SQLITE data source name")
	flag.Parse()

	ctx := context.Background()
	if err := run(ctx, os.Stdout, *addr, *dsn); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, writer io.Writer, addr string, dsn string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}
	queries := database.New(db)

	templateCache, err := web.NewTemplateCache()
	if err != nil {
		return err
	}

	server := &web.Server{
		Logger:        logger,
		Queries:       queries,
		Ctx:           ctx,
		TemplateCache: templateCache,
	}

	httpServer := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           server.Routes(),
	}

	errChan := make(chan error, 1)
	go func() {
		logger.InfoContext(ctx, "server started", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		logger.InfoContext(ctx, "server stopped")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return httpServer.Shutdown(ctx)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
