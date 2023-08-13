package main

import (
	"flag"
	"log"

	"movix/metadata/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	var grpcEnabled bool = false

	flag.BoolVar(&grpcEnabled, "grpc", false, "Enable grpc communication")
	flag.Parse()

	log.Println("starting metadata microservice")
	microservice.Start(Version, GitCommit, grpcEnabled)
	log.Println("stopping metadata microservice")
}
