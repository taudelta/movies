package microservice

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"movix/metadata/internal/config"
	"movix/metadata/internal/controller"
	"movix/metadata/internal/handler"
	"movix/metadata/internal/repository/memory"
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
		w.Write([]byte(fmt.Sprintf("metadata microservice: %s, commit: %s", version, gitCommit)))
	})

	http.HandleFunc("/metadata", handler.GetMetadata)

	http.ListenAndServe(cfg.AppAddr, nil)
}
