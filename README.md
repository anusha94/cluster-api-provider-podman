# cluster-api-provider-podman
Cluster API Provider for Podman


# Scaffolding
## use kubebuilder to start a project
```shell
kubebuilder init --domain cluster.x-k8s.io --repo github.com/anusha94/cluster-api-provider-podman
```

## create apis and controllers
```shell
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachineTemplate
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachine
kubebuilder create api --group infrastructure --version v1 --kind PodmanCluster
```

## generate manifests and implementation of runtime.Object 
Any change to the APIs, we have to generate the manifests
```shell
make manifests
make generate
```

# Implementing CAPI contract
Now that our scaffolded code is ready, we will write the provider logic. CAPI defines certain contract that each infrastructure provider has to follow. For each of our controllers, we will follow the same

## Implementing the PodmanCluster
We will follow the flow diagram mentioned in the Cluster API book for [implementing the Cluster Infrastructure](https://cluster-api.sigs.k8s.io/developer/providers/cluster-infrastructure.html#behavior)