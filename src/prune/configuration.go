package main

import (
	"encoding/json"
	"github.com/quipo/dependencysolver"
	"io/ioutil"
	"net"
	"path/filepath"
)

type Configuration struct {
	Project    string               `json:"project,omitempty"`
	Containers map[string]Container `json:"containers"`
}

type Container struct {
	// Hosts is the list of extra hosts to be added to /etc/hosts
	Hosts map[string]net.IP `json:"hosts,omitempty"`

	// CpuShares is the relative weight of CPU shares that should be given to the container
	CpuShares int32 `json:"cpu_shares,omitempty"`

	// AddCapabilities is the list of capabilities to add
	AddCapabilities []string `json:"cap_add,omitempty"`

	// DropCapabilities is the list of capabilities to drop
	DropCapabilities []string `json:"cap_drop,omitempty"`

	// CidFile is the file where we write the container ID
	CidFile string `json:"cidfile,omitempty"`

	// CpuSet is the list of CPUs on which to allow execution
	CpuSet string `json:"cpuset,omitempty"`

	// Devices is the list of host devices to be added to the container
	Devices map[string]string `json:"devices,omitempty"`

	// Dns is the list of extra DNS servers to add to the container
	DnsServers []net.IP `json:"dns,omitempty"`

	// SearchDomains is the list of search domains to add to the container
	SearchDomains []string `json:"search_domains,omitempty"`

	// EnvironmentVariables is the list of environment variables to add to the environment in the container
	EnvironmentVariables map[string]string `json:"environment,omitempty"`

	// Entrypoint is the entrypoint for the container
	Entrypoint string `json:"entrypoint,omitempty"`

	// EnvironmentFile is the file from which the environment variables can be loaded
	EnvironmentFile string `json:"environment_file,omitempty"`

	// ExposePorts is the list of ports to be exposed but not published
	ExposePorts []uint32 `json:"expose,omitempty"`

	// Hostname is the FQDN hostname to give the container
	Hostname string `json"hostname,omitempty"`

	// Ipc specifies the IPC namespace for the container
	Ipc string `json:"ipc,omitempty"`

	// Links is the list of othe containers to link to
	Links map[string]string `json:"links,omitempty"`

	// LxcConfiguration adds custom lxc options
	LxcConfiguration []string `json:"lxc,omitempty"`

	// MemoryLimit sets the limit on how much memory the container can use
	MemoryLimit string `json:"memory,omitempty"`

	// MacAddress specifies the MAC address of the container
	// TODO(zeffron: 2015-01-31) Consider using net.HardwareAddr as the type
	MacAddress string `json:"mac,omitempty"`

	// Network specifies the network of the container
	Network string `json:"network,omitempty"`

	// Ports specifies the ports to publish for the container
	Ports map[string]string `json:"ports,omitempty"`

	// Privileged specifies whether or not the container should be privileged
	Privileged bool `json:"privileged,omitempty"`

	// SecurityOptions is the list of security (such as SELinux or AppArmor) options for the container
	SecurityOptions []string `json:"security,omitempty"`

	// User specifies the username or UID to use when running the container's procees
	User string `json:"user,omitempty"`

	// Volumes is the list of volumes to bind mount
	Volumes map[string]string `json:"volumes,omitempty"`

	// VolumesFrom is the list of containers from which to mount all volumes
	VolumesFrom []string `json:"volumes_from,omitempty"`

	// WorkingDirectory is the working directory for the process run in the container
	WorkingDirectory string `json:"workdir,omitempty"`

	// Build is the path to use as the workspace when building the container
	Build string `json:"build,omitempty"`

	// Image is the image and tag used to push or fetch the container to or from the registry
	Image string `json:"image"`
}

func NewConfiguration(configurationFile string) (configuration Configuration, err error) {
	if configurationFile == "" {
		configurationFile = "prune.json"
	}

	jsonBytes, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonBytes, &configuration)
	if err != nil {
		return
	}

	for name, container := range configuration.Containers {
		if container.Build != "" && !filepath.IsAbs(container.Build) {
			var rootDir string
			rootDir, err = filepath.Abs(filepath.Dir(configurationFile))
			if err != nil {
				return
			}
			container.Build = filepath.Clean(filepath.Join(rootDir, container.Build))
			configuration.Containers[name] = container
		}
	}

	return
}

func (configuration *Configuration) OrderedContainerLayers() [][]string {
	entries := make([]dependencysolver.Entry, 0, len(configuration.Containers))
	for name, container := range configuration.Containers {
		// TODO(zeffron: 2015 01 31) Add a dependency on the container we share networking with, if there is one
		dependencies := make([]string, 0, len(container.Links)+len(container.VolumesFrom))
		for link, _ := range container.Links {
			dependencies = append(dependencies, link)
		}
		for _, volumeSource := range container.VolumesFrom {
			dependencies = append(dependencies, volumeSource)
		}
		entries = append(entries, dependencysolver.Entry{ID: name, Deps: dependencies})
	}

	return dependencysolver.LayeredTopologicalSort(entries)
}
