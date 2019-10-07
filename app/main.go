package main

import (
	"budeze/food-classifier/controller"
	"budeze/food-classifier/model"
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	f, _ := os.Create("../logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	config := model.InitConfig()
	model.InitRqlite(config.Rqlite)

	router := gin.Default()
	controller.InitController(router)

	var g errgroup.Group
	g.Go(func() error {
		return router.Run(":80")
	})

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}
	log.Println("Server exiting")
}
