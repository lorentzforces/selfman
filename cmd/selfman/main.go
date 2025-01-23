package main

import (
	"github.com/lorentzforces/selfman/internal/cli"
	"github.com/lorentzforces/selfman/internal/run"
)

func main() {
	err := cli.CreateRootCmd().Execute()
	run.FailOnErr(err)
}
