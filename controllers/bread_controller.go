/*

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
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// BreadReconciler reconciles a Bread object
type BreadReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.run-linux.com,resources=breads,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.run-linux.com,resources=breads/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
func (r *BreadReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var (
		ctx             = context.Background()
		bread           corev1alpha1.Bread
		deleteFinalizer = "onDelete"
		log             = r.Log.WithValues("bread", req.NamespacedName)
		pod             = v1.Pod{}
	)
	// The last Reconcile of deleting CR.
	if err := r.Get(ctx, req.NamespacedName, &bread); err != nil {
		log.Info(err.Error())
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// Create CR Policy
	err := r.Client.Get(ctx, req.NamespacedName, &pod)
	if errors.IsNotFound(err) {
		if err := r.OnCreate(ctx, &bread); err != nil {
			return ctrl.Result{}, err
		}
		bread.Status.Phase = pod.Status.Phase
		bread.Status.ContainerStatuses = pod.Status.ContainerStatuses
		// OnCreate() Function does not call r.Update()
		return ctrl.Result{}, r.Status().Update(ctx, &bread)
	}
	//Update CR & Pod Policy
	if pod.Spec.SchedulerName != PodSchedulingSelector(&bread) {
		log.Info("Pod: " + pod.Name + " SchedulerName changed. Ready to update the pod")
		return ctrl.Result{}, r.OnUpdate(ctx, &bread, &pod)
	}
	if !reflect.DeepEqual(pod.GetLabels(), GetPodLabel(&bread)) {
		log.Info("Pod: " + pod.Name + " label changed. Ready to update the pod")
		return ctrl.Result{}, r.OnUpdate(ctx, &bread, &pod)
	}
	// Update Pod Policy
	if pod.Status.Phase == v1.PodUnknown || pod.Status.Phase == v1.PodFailed {
		log.Info("Pod: " + pod.Name + " status is" + pod.Status.String() + " . Ready to update the pod")
		return ctrl.Result{}, r.OnUpdate(ctx, &bread, &pod)
	}
	bread.Status.Phase = pod.Status.Phase
	bread.Status.ContainerStatuses = pod.Status.ContainerStatuses
	// Delete CR Policy
	if err := r.OnDelete(ctx, req, deleteFinalizer, &bread); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *BreadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Bread{}).Watches(&source.Kind{Type: &v1.Pod{}}, &EnqueueRequest{}).
		Complete(r)
}
