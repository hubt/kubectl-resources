package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"
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

type Pod struct {
	Metadata   Metadata `yaml:"metadata"`
	Spec       PodSpec  `yaml:"spec"`
}

type PodSpec struct {
	Containers []Container
}

type Container struct {
	Name      string
	Resources Resources
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

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w,"NAMESPACE\tPOD\tCONTAINER\tCPU\tMEMORY\n")

	pods := yamlMap
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			fmt.Fprintf(w,"%s\t%s\t", pod.Metadata.Namespace, pod.Metadata.Name)
			fmt.Fprintf(w,"%s\t", container.Name)
			fmt.Fprintf(w,"%s/%s\t", container.Resources.Requests.CPU, container.Resources.Limits.CPU)
			fmt.Fprintf(w,"%s/%s\t", container.Resources.Requests.Memory, container.Resources.Limits.Memory)
			fmt.Fprintf(w,"\n")
		}
	}
	w.Flush()
}
