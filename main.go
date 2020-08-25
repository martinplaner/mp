package main

//go:generate go run generate.go assets.go

import (
	"html/template"
	"io/ioutil"

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
		zap.Reflect("config", config),
	)

	router := configureGin(config)
	router.StaticFS("/_assets/", assets)

	generator, err := GeneratorFromFile(config.File, "-")
	if err != nil {
		logger.Fatal("failed to create generator", zap.Error(err))
	}
	router.NoRoute(handlerWithGenerator(generator))

	router.Run(config.Listen)
}

func configureGin(config *Config) *gin.Engine {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	ginLogger := logger.Named("gin")
	router.Use(ginzap.RecoveryWithZap(ginLogger, true))

	t, err := loadTemplate("index.tmpl")
	if err != nil {
		logger.Fatal("failed to load templates", zap.Error(err))
	}
	router.SetHTMLTemplate(t)

	return router
}

func loadTemplate(name string) (*template.Template, error) {
	file, err := assets.Open("templates/" + name)
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
