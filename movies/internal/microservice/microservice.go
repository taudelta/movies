package microservice

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"

	"movix/movies/internal/config"
	movie "movix/movies/internal/controller"
	metadatagateway "movix/movies/internal/gateway/metadata"
	ratinggateway "movix/movies/internal/gateway/rating"
	httphandler "movix/movies/internal/handler"
	"movix/pkg/discovery"
	"movix/pkg/discovery/consul"
)

const serviceName = "metadata"

func Start(version, gitCommit string) {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	var registry *consul.Registry
	var instanceID string

	if cfg.Consul.Enabled {
		var err error
		registry, err = consul.NewRegistry(cfg.Consul.Addr)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		instanceID = discovery.GenerateInstanceID(serviceName)

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
	}

	metadataGateway := metadatagateway.New(registry, cfg.MetadataAddrs)
	ratingGateway := ratinggateway.New(registry, cfg.RatingAddrs)

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("gateway microservice: %s, commit: %s", version, gitCommit)))
	})

	http.ListenAndServe(cfg.AppAddr, nil)
}
