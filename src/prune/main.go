package main

import (
	"github.com/dynport/dgtk/cli"
)

func main() {
	router := cli.NewRouter()
	router.Register("build", &BuildRunner{}, "Build the containers.")

	err := router.RunWithArgs()
	if err != nil && err != cli.ErrorNoRoute && err != cli.ErrorHelpRequested {
		panic(err)
	}

	return
}
