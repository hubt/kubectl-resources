package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

type PodList struct {
	Items []struct {
		Metadata struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`

		Spec struct {
			Containers []struct {
				Name      string
				Resources struct {
					Limits struct {
						Memory string `yaml:"memory"`
						CPU    string `yaml:"cpu"`
					} `yaml:"limits"`
					Requests struct {
						Memory string `yaml:"memory"`
						CPU    string `yaml:"cpu"`
					} `yaml:"requests"`
					// add this to capture what kubectl top says
					Utilization struct {
						Memory string `yaml:"memory"`
						CPU    string `yaml:"cpu"`
					}
				}
			}
		} `yaml:"spec"`
	} `yaml:"items"`
}

func main() {

	var namespace string
	flag.StringVar(&namespace, "namespace", "", "Namespace")
	flag.StringVar(&namespace, "n", "", "Namespace")

	flag.Parse()

	var cmd *exec.Cmd
	if namespace == "" {
		cmd = exec.Command("kubectl", "get", "pods", "-A", "-o", "yaml")
	} else {
		cmd = exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "yaml")
	}
		
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl get pods command: %s", err)
		os.Exit(1)
	}

	var podList PodList
	err = yaml.Unmarshal(output, &podList)
	if err != nil {
		log.Fatalf("Error parsing YAML: %s", err)
	}

	parseTop(namespace,&podList)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 0, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "NAMESPACE\tPOD\tCONTAINER\tCPU:UTIL\tCPU:REQ\tCPU:LIM\tMEM:UTIL\tMEM:REQ\tMEM:LIM\n")

	pods := podList
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			fmt.Fprintf(w, "%s\t%s\t%s\t", pod.Metadata.Namespace, pod.Metadata.Name, container.Name)
			fmt.Fprintf(w, "%s\t%s\t%s\t", container.Resources.Utilization.CPU, container.Resources.Requests.CPU, container.Resources.Limits.CPU)
			fmt.Fprintf(w, "%s\t%s\t%s\n", container.Resources.Utilization.Memory, container.Resources.Requests.Memory, container.Resources.Limits.Memory)
		}
	}
	w.Flush()
}

func parseTop(namespace string, podList *PodList) error {
	var cmd *exec.Cmd
	if namespace == "" {
		cmd = exec.Command("kubectl", "top", "pods", "-A", "--containers", "--no-headers")
	} else {
		cmd = exec.Command("kubectl", "top", "pods", "-n", namespace, "--containers", "--no-headers")
	}
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl top command: %s", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		s := strings.Fields(line)
		if namespace != "" {
			s = append([]string{namespace},s...)
		}

		if len(s) != 5 {
			continue
		}

		namespace := s[0]
		pod := s[1]
		container := s[2]
		cpu := s[3]
		memory := s[4]
		for i, p := range podList.Items {
			if p.Metadata.Namespace == namespace && p.Metadata.Name == pod {
				for j, c := range p.Spec.Containers {
					if c.Name == container {
						podList.Items[i].Spec.Containers[j].Resources.Utilization.CPU = cpu
						podList.Items[i].Spec.Containers[j].Resources.Utilization.Memory = memory
						continue
					}
				}
				continue
			}
		}

	}

	return nil
}
