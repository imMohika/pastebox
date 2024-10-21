package main

import (
	"context"
	"database/sql"
	_ "embed"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"pastebox.mohika.ir/internal/database"
	"pastebox.mohika.ir/web"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	addr := flag.String("addr", ":8080", "http service address")
	dsn := flag.String("dsn", ":memory:", "SQLITE data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	queries := database.New(db)

	templateCache, err := web.NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := &web.Server{
		Logger:        logger,
		Queries:       queries,
		Ctx:           ctx,
		TemplateCache: templateCache,
	}

	logger.Info("Starting server", "addr", *addr)

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           server.Routes(),
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
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
