# cluster-api-provider-podman
Cluster API Provider for Podman


# use kubebuilder to start a project
```shell
kubebuilder init --domain cluster.x-k8s.io --repo cluster.x-k8s.io/podman
```

# create apis and controllers
```shell
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachineTemplate
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachine
kubebuilder create api --group infrastructure --version v1 --kind PodmanCluster
```

# generate manifests and implementation of runtime.Object 
Any change to the APIs, we have to generate the manifests
```shell
make manifests
make generate
```