package microservice

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"

	"movix/internal/app/metadata/config"
)

func Start(version, gitCommit string) {
	var cfg config.Config
	if err := envconfig.Parse("", &cfg); err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("metadata microservice: %s, commit: %s", version, gitCommit)))
	})

	http.ListenAndServe(cfg.AppAddr, nil)
}
