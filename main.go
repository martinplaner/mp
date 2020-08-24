package main

//go:generate go run assets_generate.go assets.go

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer logger.Sync()

	config, err := loadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	logger.Info("configuration",
		zap.String("file", config.File),
		zap.String("listen", config.Listen),
	)

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	ginLogger := logger.Named("gin")
	router.Use(ginzap.Ginzap(ginLogger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(ginLogger, true))
	router.LoadHTMLGlob("templates/*.tmpl")

	generator, err := GeneratorFromFile(config.File, "-")
	if err != nil {
		logger.Fatal("failed to create generator", zap.Error(err))
	}
	router.GET("/", handlerWithGenerator(generator))
	router.StaticFS("/_assets", assets)

	srv := &http.Server{
		Addr:    config.Listen,
		Handler: router,
	}

	go runServer(srv)
	waitForSignal()
	shutdown(srv)
}

func runServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("listen failed", zap.Error(err))
	}
}

func waitForSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func shutdown(srv *http.Server) {
	logger.Info("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}
}
