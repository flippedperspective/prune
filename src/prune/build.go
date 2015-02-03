package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
)

type BuildRunner struct {
	Verbose           bool   `cli:"opt -v --verbose desc='Enables verbose logging.'"`
	Tag               bool   `cli:"opt -t --tag desc='Tags images on successful build'"`
	ConfigurationFile string `cli:"arg desc='The prune configuration file to use'"`
}

func (runner *BuildRunner) Run() error {
	configuration, err := NewConfiguration(runner.ConfigurationFile)
	if err != nil {
		return err
	}

	color := -1
	for _, layer := range configuration.OrderedContainerLayers() {
		for _, name := range layer {
			func() {
				color += 1
				container := configuration.Containers[name]
				if container.Build == "" {
					return
				}

				// client, err := docker.NewVersionnedTLSClient("unix:///var/run/docker.sock", "", "", "", "1.16")
				client, err := docker.NewVersionedClient("unix:///var/run/docker.sock", "1.16")
				if err != nil {
					panic(err)
				}

				options := docker.BuildImageOptions{
					Name:                "",
					NoCache:             true,
					SuppressOutput:      !runner.Verbose,
					RmTmpContainer:      true,
					ForceRmTmpContainer: false,
					InputStream:         nil,
					OutputStream:        os.Stdout,
					RawJSONStream:       false,
					Remote:              "",
					Auth:                docker.AuthConfiguration{},
					AuthConfigs:         docker.AuthConfigurations{},
					ContextDir:          container.Build,
				}
				if runner.Tag {
					options.Name = container.Image
				}

				fmt.Printf("\033[38;5;%vm", color)
				fmt.Println(options)
				err = client.BuildImage(options)
				if err != nil {
					panic(err)
				}
				fmt.Print("\033[0m")
			}()
		}
	}

	return nil
}
