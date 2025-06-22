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
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	syncv1 "github.com/jacobtrvl/resonance/api/v1"
)

// ClusterSyncReconciler reconciles a ClusterSync object
type ClusterSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// Add a client for the master cluster
	MasterClient client.Client
}

// +kubebuilder:rbac:groups=sync.jacobtrvl.resonance,resources=clustersyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.jacobtrvl.resonance,resources=clustersyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sync.jacobtrvl.resonance,resources=clustersyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ClusterSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// --- ClusterSync logic (as before) ---
	agentClusterSync := &syncv1.ClusterSync{}

	// --- ReportVulnerabilities logic: watch and sync all ReportVulnerabilities to master ---
	var reportList syncv1.ReportVulnerabilitiesList

	if err := r.List(ctx, &reportList); err != nil {
		logger.Error(err, "Failed to list ReportVulnerabilities")
	} else if r.MasterClient != nil {
		for _, agentReport := range reportList.Items {
			masterReport := &syncv1.ReportVulnerabilities{}
			err := r.MasterClient.Get(ctx, client.ObjectKey{Name: agentReport.Name, Namespace: agentReport.Namespace}, masterReport)
			if err == nil {
				if !reflect.DeepEqual(agentReport.Spec.Data, masterReport.Spec.Data) {
					masterReport.Spec.Data = agentReport.Spec.Data
					if err := r.MasterClient.Update(ctx, masterReport); err != nil {
						logger.Error(err, "Failed to update ReportVulnerabilities in master cluster", "name", agentReport.Name)
					}
				}
			} else {
				logger.Error(err, "Failed to get ReportVulnerabilities in master cluster", "name", agentReport.Name)
			}
		}
	}

	/*
		// --- Deployment logic: sync Deployments from master to agent ---
		// There no watch for Deployment changes, so this will only reconcile on queing of ClusterSync due to other changes/RequeueAfter time.
		masterDeployList := &appsv1.DeploymentList{}
		masterSelector := labels.SelectorFromSet(labels.Set{"clusterSync": "true"})
		if err := r.MasterClient.List(ctx, masterDeployList, &client.ListOptions{LabelSelector: masterSelector}); err != nil {
			logger.Error(err, "Failed to list Deployments in master cluster")
		} else {
			for _, masterDeploy := range masterDeployList.Items {
				agentDeploy := &appsv1.Deployment{}
				err := r.Get(ctx, client.ObjectKey{Name: masterDeploy.Name, Namespace: masterDeploy.Namespace}, agentDeploy)
				if err != nil {
					if errors.IsNotFound(err) {
						// Create in agent if not found
						newDeploy := masterDeploy.DeepCopy()
						newDeploy.ResourceVersion = ""
						if err := r.Create(ctx, newDeploy); err != nil {
							logger.Error(err, "Failed to create Deployment in agent cluster", "name", masterDeploy.Name)
						}
					} else {
						logger.Error(err, "Failed to get Deployment in agent cluster", "name", masterDeploy.Name)
					}
				} else {
					// Update if spec differs
					if !reflect.DeepEqual(masterDeploy.Spec, agentDeploy.Spec) {
						agentDeploy.Spec = masterDeploy.Spec
						if err := r.Update(ctx, agentDeploy); err != nil {
							logger.Error(err, "Failed to update Deployment in agent cluster", "name", masterDeploy.Name)
						}
					}
				}
			}
		}*/

	// Update agentSyncStatus in ClusterSync status
	if err := r.Get(ctx, req.NamespacedName, agentClusterSync); err != nil {
		if r.MasterClient != nil {
			logger.Error(err, "Failed to get ClusterSync for status update")
		}
	} else {
		// Example: set agentSyncStatus to "Synced" and update lastSyncTime
		agentClusterSync.Status.SyncStatus = "Synced"
		if err := r.Status().Update(ctx, agentClusterSync); err != nil {
			logger.Error(err, "Failed to update ClusterSync status")
		}
	}
	fmt.Println("Reconcile called for ClusterSync:", req.NamespacedName)
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1.ClusterSync{}).
		Watches(
			&syncv1.ReportVulnerabilities{},
			&handler.EnqueueRequestForObject{},
		).
		Named("clustersync").
		Complete(r)
}
