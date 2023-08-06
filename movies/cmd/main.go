package main

import (
	"log"

	"movix/movies/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	log.Println("starting movies microservice")
	microservice.Start(Version, GitCommit)
	log.Println("stopping movies microservice")
}
