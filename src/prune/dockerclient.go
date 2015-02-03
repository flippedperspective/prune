package main

import (
	"fmt"
	"os"
)

type DockerClientWriter struct {
	Color int
}

func (client DockerClientWriter) Write(p []byte) (n int, err error) {
	b := []byte(fmt.Sprintf("\033[38;5;%vm", client.Color))
	b = append(b, p...)
	b = append(b, []byte("\033[0m")...)
	n, err = os.Stdout.Write(b)
	return
}
