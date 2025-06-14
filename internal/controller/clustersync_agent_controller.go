/*
Copyright 2025 Jacob Philip.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	syncv1 "github.com/jacobtrvl/resonance/api/v1"
)

// ClusterSyncReconciler reconciles a ClusterSync object
type ClusterSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sync.io.github.jacobtrvl,resources=clustersyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.io.github.jacobtrvl,resources=clustersyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sync.io.github.jacobtrvl,resources=clustersyncs/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ClusterSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ClusterSync instance
	clusterSync := &syncv1.ClusterSync{}
	if err := r.Get(ctx, req.NamespacedName, clusterSync); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("ClusterSync resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ClusterSync")
		return ctrl.Result{}, err
	}

	// Get the kubeconfig secret for the remote cluster
	kubeconfigSecret := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      clusterSync.Spec.RemoteClusterConfig.KubeconfigSecretName,
		Namespace: clusterSync.Spec.RemoteClusterConfig.KubeconfigSecretNamespace,
	}, kubeconfigSecret)
	if err != nil {
		log.Error(err, "Failed to get kubeconfig secret")
		clusterSync.Status.SyncStatus = "Error"
		clusterSync.Status.ErrorMessage = fmt.Sprintf("Failed to get kubeconfig secret: %v", err)
		if err := r.Status().Update(ctx, clusterSync); err != nil {
			log.Error(err, "Failed to update ClusterSync status")
		}
		return ctrl.Result{}, err
	}

	// Create a selector for resources with clusterSync=true label
	selector := labels.SelectorFromSet(labels.Set{"clusterSync": "true"})

	// List all resources with the clusterSync=true label
	// Note: This is a simplified example. In a real implementation, you would need to:
	// 1. List different types of resources (Pods, Services, ConfigMaps, etc.)
	// 2. Handle each resource type appropriately
	// 3. Implement proper error handling and retry logic
	// 4. Add proper synchronization mechanisms

	// Example for Pods
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, &client.ListOptions{LabelSelector: selector}); err != nil {
		log.Error(err, "Failed to list pods")
		clusterSync.Status.SyncStatus = "Error"
		clusterSync.Status.ErrorMessage = fmt.Sprintf("Failed to list pods: %v", err)
		if err := r.Status().Update(ctx, clusterSync); err != nil {
			log.Error(err, "Failed to update ClusterSync status")
		}
		return ctrl.Result{}, err
	}

	// TODO: Implement the actual sync logic to the remote cluster
	// This would involve:
	// 1. Creating a client for the remote cluster using the kubeconfig
	// 2. Syncing each resource to the remote cluster
	// 3. Handling conflicts and updates
	// 4. Implementing proper error handling and retry logic

	// Update the status
	clusterSync.Status.LastSyncTime = &metav1.Time{Time: time.Now()}
	clusterSync.Status.SyncStatus = "Synced"
	clusterSync.Status.ErrorMessage = ""
	if err := r.Status().Update(ctx, clusterSync); err != nil {
		log.Error(err, "Failed to update ClusterSync status")
		return ctrl.Result{}, err
	}

	// Requeue after 5 minutes to check for updates
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1.ClusterSync{}).
		Named("clustersync").
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
