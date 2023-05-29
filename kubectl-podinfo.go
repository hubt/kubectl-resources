package main

import (
	"fmt"
	"log"
	//"os"
	"os/exec"
	//"strings"

	"gopkg.in/yaml.v2"
)

type PodList struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Items      []Pod  `yaml:"items"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type Resources struct {
	Limits   Limits   `yaml:"limits"`
	Requests Requests `yaml:"requests"`
}

type Limits struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
}

type Requests struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
}

type Container struct {
	Name      string
	Resources Resources
}

type PodSpec struct {
	Containers []Container
}

type Pod struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       PodSpec  `yaml:"spec"`
}

func main() {
	// Execute "kubectl get pods -A -o  yaml" command
	cmd := exec.Command("kubectl", "get", "pods", "-A", "-o", "yaml")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl command: %s", err)
	}

	var yamlMap PodList
	err = yaml.Unmarshal(output, &yamlMap)
	if err != nil {
		log.Fatalf("Error parsing YAML: %s", err)
	}

	pods := yamlMap
	for _, pod := range pods.Items {
		fmt.Printf("Namespace: %s\n  Pod: %s\n", pod.Metadata.Namespace, pod.Metadata.Name)
		for _, container := range pod.Spec.Containers {
			fmt.Printf("    container: %s\n", container.Name)
			fmt.Printf("      CPU: %s/%s\n", container.Resources.Requests.CPU, container.Resources.Limits.CPU)
			fmt.Printf("      Memory: %s/%s\n", container.Resources.Requests.Memory, container.Resources.Limits.Memory)
		}
	}
}
