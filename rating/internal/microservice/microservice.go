package microservice

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"

	"movix/pkg/discovery"
	"movix/pkg/discovery/consul"
	"movix/rating/internal/config"
	"movix/rating/internal/controller"
	"movix/rating/internal/handler"
	"movix/rating/internal/repository/memory"
)

const serviceName = "rating"

func Start(version, gitCommit string) {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, cfg.AppAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: ", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer func() {
		if err := registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Println("Failed to deregister: ", err)
		}
	}()

	repo := memory.New()
	ctrl := controller.New(repo)
	handler := handler.NewHandler(ctrl)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("rating microservice: %s, commit: %s", version, gitCommit)))
	})

	http.HandleFunc("/rating", handler.Handle)

	http.ListenAndServe(cfg.AppAddr, nil)
}
