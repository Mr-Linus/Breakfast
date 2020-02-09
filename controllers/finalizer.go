package controllers

import (
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	"github.com/go-openapi/swag"
)

// Set the Finalizer. If the bread never has the deleteFinalizer,
// it will set the deleteFinalizer and returns true.
func (r *BreadReconciler) SetFinalizer(bread *corev1alpha1.Bread, deleteFinalizer string) bool {
	if !r.CheckFinalizer(bread, deleteFinalizer) {
		bread.ObjectMeta.Finalizers = append(bread.ObjectMeta.Finalizers, deleteFinalizer)
		return true
	}
	return false
}

// Check the Bread has the deleteFinalizer.
func (r *BreadReconciler) CheckFinalizer(bread *corev1alpha1.Bread, deleteFinalizer string) bool {
	return swag.ContainsStrings(bread.ObjectMeta.Finalizers, deleteFinalizer)
}

// Delete the deleteFinalizer.
func (r *BreadReconciler) DeleteFinalizer(bread *corev1alpha1.Bread, deleteFinalizer string) {
	if r.CheckFinalizer(bread, deleteFinalizer) {
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
