/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrastructurev1 "github.com/anusha94/cluster-api-provider-podman/api/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// PodmanMachineReconciler reconciles a PodmanMachine object
type PodmanMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanmachines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodmanMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PodmanMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile request received")

	// Fetch the PodmanMachien instance
	podmanMachine := &infrastructurev1.PodmanMachine{}
	err := r.Client.Get(ctx, req.NamespacedName, podmanMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Machine
	machine, err := util.GetOwnerMachine(ctx, r.Client, podmanMachine.ObjectMeta)
	if err != nil {
		logger.Error(err, "failed to get Owner Machine")
		return ctrl.Result{}, err
	}

	if machine == nil {
		logger.Info("Waiting for Machine Controller to set OwnerRef on PodmanMachine")
		return ctrl.Result{}, nil
	}

	// Fetch the Cluster
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, podmanMachine.ObjectMeta)
	if err != nil {
		logger.Error(err, "PodmanMachine owner Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, err
	}

	if cluster == nil {
		logger.Info(fmt.Sprintf("Please associate this machine with a cluster using the label %s: <name of cluster>", clusterv1.ClusterLabelName))
		return ctrl.Result{}, nil
	}
	logger = logger.WithValues("cluster", cluster.Name)

	podmanCluster := &infrastructurev1.PodmanCluster{}
	infraClusterName := client.ObjectKey{
		Namespace: podmanMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err = r.Client.Get(ctx, infraClusterName, podmanCluster); err != nil {
		logger.Error(err, "failed to get infra cluster")
		return ctrl.Result{}, nil
	}

	helper, _ := patch.NewHelper(podmanMachine, r.Client)
	defer func() {
		if err = helper.Patch(ctx, podmanMachine); err != nil && reterr == nil {
			logger.Error(err, "failed to patch podmanMachine")
			reterr = err
		}
	}()

	// Handle deleted machines
	if !podmanMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, podmanMachine)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, podmanMachine)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodmanMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1.PodmanMachine{}).
		Complete(r)
}

func (r *PodmanMachineReconciler) reconcileDelete(ctx context.Context, podmanMachine *infrastructurev1.PodmanMachine) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Deleting PodmanMachine")

	controllerutil.RemoveFinalizer(podmanMachine, infrastructurev1.MachineFinalizer)
	return reconcile.Result{}, nil
}

func (r *PodmanMachineReconciler) reconcileNormal(ctx context.Context, podmanMachine *infrastructurev1.PodmanMachine) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling PodmanMachine")

	// TODO: Write reconcile logic
	return reconcile.Result{}, nil
}
