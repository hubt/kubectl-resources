package main

import (
	"fmt"
	"log"
	//"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type ContainerInfo struct {
	Name        string
	CPURequests string
	CPULimits   string
	MemRequests string
	MemLimits   string
}

type PodInfo struct {
	Namespace   string
	Name        string
	Containers  []ContainerInfo
	CPURequests string
	CPULimits   string
	MemRequests string
	MemLimits   string
}

func main() {
	// Execute "kubectl get namespaces" command
	cmd := exec.Command("kubectl", "get", "namespaces", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl command: %s", err)
	}

	// Parse the output and get the list of namespaces
	namespaces := parseNamespaces(string(output))

	// Iterate over each namespace and get the pod information
	var podInfos []PodInfo
	for _, namespace := range namespaces {
		pods := getPods(namespace)
		for _, pod := range pods {
			podInfo := getPodInfo(namespace, pod)
			podInfos = append(podInfos, podInfo)
		}
	}

	// Display the pod information
	for _, podInfo := range podInfos {
		fmt.Printf("Namespace: %s, Pod: %s\n", podInfo.Namespace, podInfo.Name)
		fmt.Printf("  CPU Requests: %s\n", podInfo.CPURequests)
		fmt.Printf("  CPU Limits: %s\n", podInfo.CPULimits)
		fmt.Printf("  Memory Requests: %s\n", podInfo.MemRequests)
		fmt.Printf("  Memory Limits: %s\n\n", podInfo.MemLimits)
	}
}

func parseNamespaces(output string) []string {
	lines := strings.Split(output, "\n")
	var namespaces []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			namespaces = append(namespaces, fields[0])
		}
	}
	return namespaces
}

func getPods(namespace string) []string {
	// Execute "kubectl get pods" command with the given namespace
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "--no-headers", "-o", "custom-columns=:.metadata.name")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl command: %s", err)
	}

	// Parse the output and get the list of pod names
	lines := strings.Split(string(output), "\n")
	var pods []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			pods = append(pods, line)
		}
	}
	return pods
}

func getPodInfo(namespace, pod string) PodInfo {
	// Execute "kubectl get pod" command with the given namespace and pod name in YAML format
	cmd := exec.Command("kubectl", "get", "pod", pod, "-n", namespace, "-o", "yaml")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing kubectl command: %s", err)
	}

	// Parse the YAML output and get the CPU and memory requests and limits
	//cpuRequests := parseYAMLValue(output, "cpu-requests")
	//cpuLimits := parseYAMLValue(output, "cpu-limits")
	//memRequests := parseYAMLValue(output, "memory-requests")
	//memLimits := parseYAMLValue(output, "memory-limits")

	var podInfo = PodInfo{
		Namespace: namespace,
		Name: pod,
	}
	parseYAMLValue(output,&podInfo)
	fmt.Printf("output: %s\n",podInfo)
	//fmt.Printf("cpuRequests: %s\n",cpuRequests)
	//fmt.Printf("cpuLimits: %s\n",cpuLimits)
	//fmt.Printf("memRequests: %s\n",memRequests)
	//fmt.Printf("memLimits: %s\n",memLimits)
	//return string(podInfo)
	return podInfo
}


func parseYAMLValue(yamlData []byte,podInfo *PodInfo) error {
	var yamlMap map[string]interface{}
	err := yaml.Unmarshal(yamlData, &yamlMap)
	if err != nil {
		log.Fatalf("Error parsing YAML: %s", err)
	}
	//fmt.Printf("yaml: %s\n",yamlMap)
	//fmt.Printf("\n\nyaml1: %s\n", yamlMap["spec"].(map[string]interface{})["containers"])

	//containers, ok := yamlMap["spec"].(map[string]interface{})["containers"].([]interface{})
	//containers, ok := yamlMap["spec"].(map[string]interface{})["containers"].([]interface{})
	containers, ok := yamlMap["spec"].(map[interface{}]interface{})["containers"].([]interface{})

	if !ok {
		return nil
	}

	for _, container := range containers {
		var containerInfo ContainerInfo
		podInfo.Containers = append(podInfo.Containers,containerInfo)
		containerData := container.(map[interface{}]interface{})
		if name, ok := containerData["name"].(string); ok {
			containerInfo.Name = name
		}
			
		resources, ok := containerData["resources"].(map[interface{}]interface{})
		if ok {
			limits, ok := resources["limits"].(map[interface{}]interface{})
			if ok {
				memLimits, ok := limits["memory"].(string)
				if ok {
					containerInfo.MemLimits = memLimits
					podInfo.MemLimits += memLimits
				}
				cpuLimits, ok := limits["cpu"].(string)
				if ok {
					containerInfo.CPULimits = cpuLimits
					podInfo.CPULimits += cpuLimits
				}
			}

			requests, ok := resources["requests"].(map[interface{}]interface{})
			if ok {
				memRequests, ok := requests["memory"].(string)
				if ok {
					containerInfo.MemRequests = memRequests
					podInfo.MemRequests += memRequests
				}
				cpuRequests, ok := requests["cpu"].(string)
				if ok {
					containerInfo.CPURequests = cpuRequests
					podInfo.CPURequests += cpuRequests
				}
			}
		}
	}

	return nil
}

