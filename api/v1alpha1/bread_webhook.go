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

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var breadlog = logf.Log.WithName("bread-resource")

func (r *Bread) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-core-run-linux-com-v1alpha1-bread,mutating=true,failurePolicy=fail,groups=core.run-linux.com,resources=breads,verbs=create;update,versions=v1alpha1,name=mbread.kb.io

var _ webhook.Defaulter = &Bread{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Bread) Default() {
	breadlog.Info("default", "name", r.Name)
	if r.Spec.Task.Type == "ssh" {
		r.Spec.Task.Command = ""
	}
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-core-run-linux-com-v1alpha1-bread,mutating=false,failurePolicy=fail,groups=core.run-linux.com,resources=breads,versions=v1alpha1,name=vbread.kb.io

var _ webhook.Validator = &Bread{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Bread) ValidateCreate() error {

	breadlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.ValidateBread()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Bread) ValidateUpdate(old runtime.Object) error {
	breadlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.ValidateBread()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Bread) ValidateDelete() error {
	breadlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Bread) ValidateTask() *field.Error {
	if r.Spec.Task.Type != "ssh" && r.Spec.Task.Type != "train" {
		return field.Invalid(field.NewPath("spec").Child("task").Child("type"),
			r.Spec.Task.Type,
			"Task Type can't set to:"+
				r.Spec.Task.Type)
	}
	if r.Spec.Task.Type == "train" && r.Spec.Task.Command == "" {
		return field.Invalid(field.NewPath("spec").Child("task").Child("command"),
			r.Spec.Task.Type,
			"Task Command can't set to Null")
	}
	return nil
}

func (r *Bread) ValidateFreamwork() *field.Error {
	if r.Spec.Framework.Name != "tensorflow" && r.Spec.Framework.Name != "pytorch" && r.Spec.Framework.Name != "cuda" {
		return field.Invalid(
			field.NewPath("spec").Child("framework").Child("name"),
			r.Spec.Framework.Name,
			"Framework can't set to:"+
				r.Spec.Framework.Name)
	}
	return nil
}

func (r *Bread) ValidateBread() error {
	var allErrs field.ErrorList
	if err := r.ValidateFreamwork(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.ValidateTask(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "core.run-linux.com", Kind: "Bread"},
		r.Name, allErrs)
}
