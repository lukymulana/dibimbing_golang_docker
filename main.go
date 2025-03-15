package main

import (
	"boilerplate/internal/config"
	helloworld "boilerplate/internal/hello-world"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	var logger *zap.Logger
	var mode string

	switch cfg.Env {
	case "prod":
		mode = gin.ReleaseMode
		l, _ := zap.NewProduction()
		logger = l
	default:
		mode = gin.DebugMode
		l, _ := zap.NewDevelopment()
		logger = l
	}

	gin.SetMode(mode)

	r := gin.New()
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
	}))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors.Default())

	// Init Hello World Router
	helloWorldHandler := helloworld.NewHandler()
	helloWorldRouter := helloworld.NewRouter(helloWorldHandler, r.RouterGroup)
	helloWorldRouter.Register()

	r.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")
}
