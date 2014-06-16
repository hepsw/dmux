package main

import (
	"github.com/fsouza/go-dockerclient"
	"os"
	"os/exec"
	"strings"
)

func PauseContainer(containerName string) {
	cmd := exec.Command("docker", "pause", containerName)
	RunCommand(cmd)
}

func UnpauseContainer(containerName string) {
	cmd := exec.Command("docker", "unpause", containerName)
	RunCommand(cmd)
}

func DeleteContainer(containerName string) {
	cmd := exec.Command("docker", "rm", "-f", containerName)
	RunCommand(cmd)
}

func SaveContainer(containerName string, imageName string) {
	dockerAddr := os.Getenv("DOCKER_HOST")
	client, err := docker.NewClient(dockerAddr)
	assert(err)

	commitOption := docker.CommitContainerOptions{
		Container:  containerName,
		Repository: imageName,
	}

	image, err := client.CommitContainer(commitOption)
	assert(err)
	debug(image.ID)
}

func HasContainer(containerName string, fromAll bool) bool {

	dockerAddr := os.Getenv("DOCKER_HOST")
	client, err := docker.NewClient(dockerAddr)
	assert(err)

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: fromAll,
	})
	assert(err)

	for _, listing := range containers {
		debug(listing.Names[0])
		workContainerName := strings.TrimLeft(listing.Names[0], "/")
		if workContainerName == containerName {
			return true
		}
	}

	return false
}

func IsExited(containerName string) bool {

	dockerAddr := os.Getenv("DOCKER_HOST")
	client, err := docker.NewClient(dockerAddr)
	assert(err)

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	assert(err)

	for _, listing := range containers {
		debug(listing.Names[0])
		debug(listing.Status)

		workContainerName := strings.TrimLeft(listing.Names[0], "/")
		if workContainerName != containerName {
			continue
		}

		containerStatus := listing.Status
		if strings.HasPrefix(containerStatus, "Exited") {
			return true
		}
	}
	return false
}

func IsPaused(containerName string) bool {

	dockerAddr := os.Getenv("DOCKER_HOST")
	client, err := docker.NewClient(dockerAddr)
	assert(err)

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: false,
	})
	assert(err)

	for _, listing := range containers {
		debug(listing.Names[0])
		debug(listing.Status)

		workContainerName := strings.TrimLeft(listing.Names[0], "/")
		if workContainerName != containerName {
			continue
		}

		containerStatus := listing.Status
		if strings.HasSuffix(containerStatus, "(Paused)") {
			return true
		}
	}
	return false
}
