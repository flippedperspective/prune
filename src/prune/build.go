package main

import (
	"github.com/fsouza/go-dockerclient"
	"sync"
)

type BuildRunner struct {
	Verbose           bool   `cli:"opt -v --verbose desc='Enables verbose logging.'"`
	Tag               bool   `cli:"opt -t --tag desc='Tags images on successful build'"`
	NoCache           bool   `cli:"opt --no-cache desc='Build without using the docker cache.'"`
	ConfigurationFile string `cli:"arg desc='The prune configuration file to use'"`
}

func (runner *BuildRunner) Run() error {
	configuration, err := NewConfiguration(runner.ConfigurationFile)
	if err != nil {
		return err
	}

	color := 0
	var wg sync.WaitGroup
	for _, layer := range configuration.OrderedContainerLayers() {
		for _, name := range layer {
			wg.Add(1)
			go func() {
				defer wg.Done()
				container := configuration.Containers[name]
				if container.Build == "" {
					return
				}

				color += 1

				client, err := docker.NewVersionedClient("unix:///var/run/docker.sock", "1.16")
				if err != nil {
					panic(err)
				}

				options := docker.BuildImageOptions{
					Name:                "",
					NoCache:             runner.NoCache,
					SuppressOutput:      !runner.Verbose,
					RmTmpContainer:      true,
					ForceRmTmpContainer: false,
					InputStream:         nil,
					OutputStream:        DockerClientWriter{Color: color},
					RawJSONStream:       false,
					Remote:              "",
					Auth:                docker.AuthConfiguration{},
					AuthConfigs:         docker.AuthConfigurations{},
					ContextDir:          container.Build,
				}
				if runner.Tag {
					options.Name = container.Image
				}

				err = client.BuildImage(options)
				if err != nil {
					panic(err)
				}
			}()
		}
	}
	wg.Wait()

	return nil
}
