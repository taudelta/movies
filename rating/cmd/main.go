package main

import (
	"movix/rating/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	microservice.Start(Version, GitCommit)
}
