package main

import (
	"movix/metadata/internal/microservice"
)

var Version string = "undefined"
var GitCommit string = "undefined"

func main() {
	microservice.Start(Version, GitCommit)
}
