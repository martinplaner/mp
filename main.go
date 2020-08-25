package main

//go:generate go run generate.go assets.go templates.go

import (
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	logger.Info("effective configuration",
		zap.Reflect("config", config),
	)

	router := configureGin(config)

	generator, err := GeneratorFromFile(config.File, "-")
	if err != nil {
		logger.Fatal("failed to create generator", zap.Error(err))
	}
	router.GET("/", handlerWithGenerator(generator))
	router.StaticFS("/_assets/", assets)

	srv := &http.Server{
		// TODO handle timeouts?
		Addr:    config.Listen,
		Handler: router,
	}

	runServer(srv)
}

func configureGin(config *Config) *gin.Engine {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	t, err := loadTemplate("index.tmpl")
	if err != nil {
		logger.Fatal("failed to load templates", zap.Error(err))
	}
	router.SetHTMLTemplate(t)

	return router
}

func loadTemplate(name string) (*template.Template, error) {
	file, err := templates.Open(name)
	if err != nil {
		logger.Fatal("failed to load template", zap.String("name", name), zap.Error(err))
	}
	defer file.Close()

	h, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	t, err := template.New(name).Parse(string(h))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func runServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}
}
