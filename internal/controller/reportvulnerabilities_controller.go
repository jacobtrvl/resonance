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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	syncv1 "github.com/jacobtrvl/resonance/api/v1"
)

// ReportVulnerabilitiesReconciler reconciles ReportVulnerabilities custom resources
// +kubebuilder:rbac:groups=sync.jacobtrvl.resonance,resources=reportvulnerabilities,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.jacobtrvl.resonance,resources=reportvulnerabilities/status,verbs=get;update;patch

type ReportVulnerabilitiesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ReportVulnerabilitiesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ReportVulnerabilities instance that triggered the reconcile
	report := &syncv1.ReportVulnerabilities{}
	if err := r.Get(ctx, req.NamespacedName, report); err != nil {
		logger.Error(err, "Failed to get ReportVulnerabilities")
		return ctrl.Result{}, err
	}

	// Log the update event
	logger.Info("ReportVulnerabilities updated", "name", report.Name, "namespace", report.Namespace, "data", report.Spec.Data)

	return ctrl.Result{}, nil
}

func (r *ReportVulnerabilitiesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1.ReportVulnerabilities{}).
		Named("reportvulnerabilities").
		Complete(r)
}
