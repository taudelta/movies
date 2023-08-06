package main

import (
	"log"

	"movix/rating/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	log.Println("starting rating microservice")
	microservice.Start(Version, GitCommit)
	log.Println("stopping rating microservice")
}
