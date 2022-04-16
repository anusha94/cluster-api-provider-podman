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

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrastructurev1 "github.com/anusha94/cluster-api-provider-podman/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// PodmanClusterReconciler reconciles a PodmanCluster object
type PodmanClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=podmanclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodmanCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PodmanClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the PodmanCluster instance
	podmanCluster := &infrastructurev1.PodmanCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, podmanCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, podmanCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Waiting for Cluster Controller to set OwnerRef on PodmanCluster")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(podmanCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the PodmanCluster object and status after each reconciliation.
	defer func() {
		if err := patchPodmanCluster(ctx, patchHelper, podmanCluster); err != nil {
			log.Error(err, "failed to patch PodmanCluster")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !controllerutil.ContainsFinalizer(podmanCluster, infrastructurev1.ClusterFinalizer) {
		controllerutil.AddFinalizer(podmanCluster, infrastructurev1.ClusterFinalizer)
		return ctrl.Result{}, nil
	}

	// Handle deleted clusters
	if !podmanCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, podmanCluster)
	}

	// Handle non-deleted clusters
	return r.reconcileNormal(ctx, podmanCluster)
}

func (r PodmanClusterReconciler) reconcileDelete(ctx context.Context, byoCluster *infrastructurev1.PodmanCluster) (reconcile.Result, error) {
	_ = log.FromContext(ctx)

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(byoCluster, infrastructurev1.ClusterFinalizer)

	return ctrl.Result{}, nil
}

func (r PodmanClusterReconciler) reconcileNormal(ctx context.Context, podmanCluster *infrastructurev1.PodmanCluster) (reconcile.Result, error) {
	// If the ByoCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(podmanCluster, infrastructurev1.ClusterFinalizer)

	podmanCluster.Status.Ready = true

	return reconcile.Result{}, nil
}

func patchPodmanCluster(ctx context.Context, patchHelper *patch.Helper, podmanCluster *infrastructurev1.PodmanCluster) error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding it during the deletion process).
	conditions.SetSummary(podmanCluster,
		conditions.WithStepCounterIf(podmanCluster.ObjectMeta.DeletionTimestamp.IsZero()),
	)

	// Patch the object, ignoring conflicts on the conditions owned by this controller.
	return patchHelper.Patch(
		ctx,
		podmanCluster,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
		}},
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodmanClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1.PodmanCluster{}).
		Complete(r)
}
