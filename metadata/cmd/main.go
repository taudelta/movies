package main

import (
	"log"

	"movix/metadata/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	log.Println("starting metadata microservice")
	microservice.Start(Version, GitCommit)
	log.Println("stopping metadata microservice")
}
