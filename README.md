
this kubectl plugin queries a kubernetes cluster for running pods, parses out the container requests and limits for memory and cpu.
then does a kubectl top pods and shows the utilization, requests, and limits together for easier reference

## Build:
either `make` or `go build *.go` 

## Run:

```
./kubectl-resources
```
can be installed as any other kubectl plugin, so then you can just run:
```
kubectl resources
```


```
NAMESPACE  |POD                                    |CONTAINER             |CPU:UTIL|CPU:REQ|CPU:LIM|MEM:UTIL|MEM:REQ|MEM:LIM
kube-system|helm-install-traefik-crd-zkw2n         |helm                  |        |       |       |        |       |
kube-system|helm-install-traefik-tn7z9             |helm                  |        |       |       |        |       |
kube-system|traefik-7cd4fcff68-nq8b8               |traefik               |2m      |       |       |68Mi    |       |
kube-system|coredns-b96499967-b6fwc                |coredns               |2m      |100m   |       |36Mi    |70Mi   |170Mi
kube-system|local-path-provisioner-7b7dc8d6f5-76s4n|local-path-provisioner|1m      |       |       |22Mi    |       |
kube-system|metrics-server-668d979685-przfc        |metrics-server        |4m      |100m   |       |43Mi    |70Mi   |
kube-system|svclb-traefik-cc17e230-kpdp8           |lb-tcp-80             |0m      |       |       |1Mi     |       |
kube-system|svclb-traefik-cc17e230-kpdp8           |lb-tcp-443            |0m      |       |       |0Mi     |       |
default    |nginx-8cb56b9b9-f284m                  |nginx                 |0m      |300m   |3      |13Mi    |300M   |3G
default    |nginx-8cb56b9b9-f284m                  |nginx2                |        |200m   |2      |        |200M   |2G
```
