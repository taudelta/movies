package main

import (
	"movix/internal/movies/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	microservice.Start(Version, GitCommit)
}
