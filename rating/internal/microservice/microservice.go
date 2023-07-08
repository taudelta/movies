package microservice

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"movix/rating/internal/config"
	"movix/rating/internal/controller"
	"movix/rating/internal/handler"
	"movix/rating/internal/repository/memory"
)

func Start(version, gitCommit string) {
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	repo := memory.New()

	ctrl := controller.New(repo)

	handler := handler.NewHandler(ctrl)

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("rating microservice: %s, commit: %s", version, gitCommit)))
	})

	http.HandleFunc("/rating", handler.Handle)

	http.ListenAndServe(cfg.AppAddr, nil)
}
