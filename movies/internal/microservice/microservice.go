package microservice

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"movix/movies/internal/config"
	movie "movix/movies/internal/controller"
	metadatagateway "movix/movies/internal/gateway/metadata"
	ratinggateway "movix/movies/internal/gateway/rating"
	httphandler "movix/movies/internal/handler"
)

func Start(version, gitCommit string) {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("gateway microservice: %s, commit: %s", version, gitCommit)))
	})

	http.ListenAndServe(cfg.AppAddr, nil)
}
