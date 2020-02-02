package controllers

import (
	"context"
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	"github.com/go-openapi/swag"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OnCreate is used to create pod. It will not update the CR.
func (r *BreadReconciler) OnCreate(ctx context.Context, bread *corev1alpha1.Bread) error {
	if TaskIsSSH(bread) {
		return r.CreateSSHPod(ctx, bread)
	}
	return r.CreateTaskPod(ctx, bread)
}

// OnDelete is used to delete pod. It will be update the CR.
// OnDelete will judge whether CR needs to be deleted.
// If the CR needs to be deleted, OnDelete will delete the pod
// by deleteFinalizer.
func (r *BreadReconciler) OnDelete(ctx context.Context, deleteFinalizer string, bread *corev1alpha1.Bread) error {
	if bread.ObjectMeta.DeletionTimestamp.IsZero() {
		if !swag.ContainsStrings(bread.ObjectMeta.Finalizers, deleteFinalizer) {
			bread.ObjectMeta.Finalizers = append(bread.ObjectMeta.Finalizers, deleteFinalizer)
		}
	} else {
		if swag.ContainsStrings(bread.ObjectMeta.Finalizers, deleteFinalizer) {
			bList := &corev1alpha1.BreadList{}
			err := r.List(ctx, bList, client.InNamespace(bread.Namespace))
			if err != nil {
				return err
			}
			for _, b := range bList.Items {
				pod := v1.Pod{}
				err := r.Client.Get(ctx, client.ObjectKey{Namespace: b.Namespace, Name: b.Name}, &pod)
				if err != nil {
					r.Log.Info(err.Error())
				} else {
					if err = r.Delete(ctx, &pod); err != nil {
						r.Log.Info(err.Error())
					}
				}
				err = r.Delete(ctx, &b)
				if err != nil {
					return err
				}
			}
			bread.ObjectMeta.Finalizers =
				func(value string, chain []string) []string {
					var list = chain
					for i, v := range chain {
						if v == value {
							list = append(list[:i], list[i+1:]...)
						}
					}
					return list
				}(deleteFinalizer, bread.ObjectMeta.Finalizers)
		}
	}
	return r.Update(ctx, bread)
}

// OnUpdate is used to update pod. It will not be update the CR.
func (r *BreadReconciler) OnUpdate(ctx context.Context, bread *corev1alpha1.Bread, pod *v1.Pod) error {
	if err := r.Client.Delete(ctx, pod); err != nil {
		return err
	}
	return nil
}
