# cluster-api-provider-podman
Cluster API Provider for Podman


# use kubebuilder to start a project
kubebuilder init --domain cluster.x-k8s.io --repo cluster.x-k8s.io/podman

# create apis and controllers
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachineTemplate
kubebuilder create api --group infrastructure --version v1 --kind PodmanMachine
kubebuilder create api --group infrastructure --version v1 --kind PodmanCluster