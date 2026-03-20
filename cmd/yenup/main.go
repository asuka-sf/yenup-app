package main

import (
	"context"
	"fmt"
	"log"

	"yenup/internal/config"
	"yenup/internal/registry"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	// load config from config.go
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//  GCSClient
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer gcsClient.Close()

	// create registry from registry.go
	reg, err := registry.NewRegistry(cfg, gcsClient)
	if err != nil {
		log.Fatal(err)
	}

	// create app handler from registry
	appHandler := reg.AppHandler

	// create router from gin
	r := gin.Default()

	// register routes from route.go
	appHandler.RegisterRoutes(r)

	// run server
	r.Run(fmt.Sprintf(":%s", cfg.AppPort))

}
