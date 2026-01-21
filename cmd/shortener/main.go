package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/mabishka/lupanova/internal/auth"
	"github.com/mabishka/lupanova/internal/compress"
	"github.com/mabishka/lupanova/internal/config"
	"github.com/mabishka/lupanova/internal/handler"
	"github.com/mabishka/lupanova/internal/logger"
	"github.com/mabishka/lupanova/internal/model"
	"github.com/mabishka/lupanova/internal/repository/audit"
	"github.com/mabishka/lupanova/internal/repository/connloader"
	"github.com/mabishka/lupanova/internal/repository/fileloader"
)

const stopTimeout = 5 * time.Second

func main() {
	if err := new(context.WithCancelCause(context.Background())); err != nil {
		log.Fatalf("exist with error: %v", err)
	}
}

func new(ctx context.Context, fnCancel context.CancelCauseFunc) error {

	config := config.New()
	if err := logger.InitLogger(config.GetLogLevel()); err != nil {
		panic(err)
	}

	logger.Log().Info("config",
		zap.String("serverAddress", config.GetServerAddress()),
		zap.String("baseAddress", config.GetBaseAddress()),
		zap.String("logLevel", config.GetLogLevel()),
		zap.String("fileName", config.GetFileName()),
		zap.String("connAddress", config.GetConnAddress()),
	)

	server := handler.New(config.GetBaseAddress())

	var loader model.StorageLoader
	var conn model.ConnLoader
	if config.GetConnAddress() != "" {
		loader = connloader.New(config.GetConnAddress())
		conn, _ = loader.(model.ConnLoader)
		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Info("conn not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("conn storage usage")
		}

	}

	if loader == nil && config.GetFileName() != "" {
		loader = fileloader.New(config.GetFileName())

		if err := server.Load(context.Background(), loader); err != nil {
			logger.Log().Error("file not loaded", zap.Error(err))
			loader = nil
		} else {
			logger.Log().Info("file storage usage")
		}
	}

	if loader == nil {
		logger.Log().Info("memory storage usage")
	}

	connServer := handler.NewConn(conn)

	auditEvent := audit.NewAuditEvent()

	if config.GetAuditFile() != "" {
		auditEvent.Register(audit.NewFileObserver(config.GetAuditFile()))
	}

	if config.GetAuditAddress() != "" {
		auditEvent.Register(audit.NewFileObserver(config.GetAuditFile()))
	}

	server.SetAudit(auditEvent)

	router := chi.NewRouter()

	router.Use(logger.WithLogging)
	router.Use(compress.WithCompress)
	router.Use(auth.WithAuth)

	router.Mount("/debug", middleware.Profiler())

	router.Post("/", server.HandlerPostFull)
	router.Post("/api/shorten", server.HandlerPostFullJSON)
	router.Post("/api/shorten/batch", server.HandlerPostBatch)
	router.Get("/{id}", server.HandlerGetFull)
	router.Get("/ping", connServer.HandlerGetPing)
	router.Delete("/api/user/urls", server.HandlerDelete)
	router.Get("/api/user/urls", server.HandlerGetUser)

	if err := run(ctx, &http.Server{
		Addr:         config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}); err != nil {
		fnCancel(err)
		return err
	}
	fnCancel(nil)

	return nil
}

func run(ctx context.Context, srv *http.Server) error {

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case s := <-sigint:
			logger.Log().Info("stop with signal", zap.String("signal", s.String()))
		case <-ctx.Done():
			logger.Log().Info("stop with context", zap.Error(context.Cause(ctx)))
		}

		stopCtx, cancel := context.WithTimeoutCause(context.Background(), stopTimeout, fmt.Errorf("server Shutdown with timeout %v", stopTimeout))
		defer cancel()
		if err := srv.Shutdown(stopCtx); err != nil {
			logger.Log().Info("HTTP server shutdown", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log().Info("HTTP server ListenAndServe", zap.Error(err))
		return err
	}

	logger.Log().Info("exit")
	return nil
}
