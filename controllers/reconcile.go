package controllers

import (
	"context"
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// OnCreate is used to create pod. It will not update the CR.
func (r *BreadReconciler) OnCreate(ctx context.Context, bread *corev1alpha1.Bread) error {
	log := r.Log.WithName("Create")
	if TaskIsSSH(bread) {
		if labels := bread.GetLabels(); labels["kubernetes.io/hostname"] != "" {
			log.Info("Create SSH Pod" + bread.Name + "with selected node: " + labels["kubernetes.io/hostname"])
			return r.CreateSSHPodWithNodeSelected(ctx, bread, labels)
		} else {
			log.Info("Create SSH Pod: " + bread.Name)
			return r.CreateSSHPodWithoutNodeSelected(ctx, bread)
		}
	} else {
		log.Info("Create Task Pod: " + bread.Name)
		return r.CreateTaskPod(ctx, bread)
	}
}

// OnDelete is used to delete pod. It will be update the CR.
// OnDelete will judge whether CR needs to be deleted.
// If the CR needs to be deleted, OnDelete will return true and
// delete Finalizer.
func (r *BreadReconciler) OnDelete(ctx context.Context, req ctrl.Request, deleteFinalizer string, bread *corev1alpha1.Bread) error {
	log := r.Log.WithName("Delete")
	if !bread.ObjectMeta.DeletionTimestamp.IsZero() {
		if r.CheckFinalizer(bread, deleteFinalizer) {
			log.Info("Delete Pod " + bread.Name)
			if err := r.DeletePod(ctx, req); err != nil && !errors.IsNotFound(err) {
				return err
			}
			r.DeleteFinalizer(bread, deleteFinalizer)
		}
	}
	return nil
}

// Check the Bread is need to be deleted
func (r *BreadReconciler) NeedToDelete(bread *corev1alpha1.Bread) bool {
	return !bread.ObjectMeta.DeletionTimestamp.IsZero()
}

// OnUpdate is used to update pod. It will not be update the CR.
func (r *BreadReconciler) OnUpdate(ctx context.Context, pod *v1.Pod) error {
	log := r.Log.WithName("Update")
	log.Info("Update Pod " + pod.Name)
	if err := r.Client.Delete(ctx, pod); err != nil {
		return err
	}
	return nil
}
