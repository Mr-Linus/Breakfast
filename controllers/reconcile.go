package controllers

import (
	"context"
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	"github.com/go-openapi/swag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *BreadReconciler) OnCreate(ctx context.Context, bread *corev1alpha1.Bread) error {
	if TaskIsSSH(bread) {
		return r.CreateSSHPod(ctx, bread)
	} else {
		return nil
	}
}

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
			bread.ObjectMeta.Finalizers = DeleteElem(deleteFinalizer, bread.ObjectMeta.Finalizers)
		}
	}
	return r.Update(ctx, bread)
}

func DeleteElem(value string, chain []string) []string {
	var list = chain
	for i, v := range chain {
		if v == value {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func GetTaskType(bread *corev1alpha1.Bread) string {
	return bread.Spec.Task.Type
}

func GetPodLabel(bread *corev1alpha1.Bread) map[string]string {
	labels := map[string]string{"bread": bread.Name}

	if bread.Spec.Scv.Level != "" {
		labels["scv/Label"] = bread.Spec.Scv.Level
	}
	if bread.Spec.Scv.Gpu != "0" {
		labels["scv/Gpu"] = bread.Spec.Scv.Gpu
	}
	if bread.Spec.Scv.Memory != "" {
		labels["scv/FreeMemory"] = bread.Spec.Scv.Memory
	}
	return labels
}

func TaskIsSSH(bread *corev1alpha1.Bread) bool {
	return GetTaskType(bread) == "ssh"
}

func GetPodImage(bread *corev1alpha1.Bread) string {
	if bread.Spec.Scv.Gpu != "0" {
		return "registry.cn-hangzhou.aliyuncs.com/njupt-isl/" +
			bread.Spec.Framework.Name +
			"-gpu:" +
			bread.Spec.Framework.Version
	}
	return "registry.cn-hangzhou.aliyuncs.com/njupt-isl/" +
		bread.Spec.Framework.Name +
		"-cpu:" +
		bread.Spec.Framework.Version
}

func (r *BreadReconciler) CreateSSHPod(ctx context.Context, bread *corev1alpha1.Bread) error {
	var sshPod = v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: bread.Namespace,
			Name:      bread.Name,
			Labels:    GetPodLabel(bread),
		},
		Spec: v1.PodSpec{
			//SchedulerName:"yoda-scheduler"
			Containers: []v1.Container{
				{
					Name:  bread.Name,
					Image: GetPodImage(bread),
					Env: []v1.EnvVar{
						{
							Name:  "PYTHONUNBUFFERED",
							Value: "0",
						},
					},

					Resources: v1.ResourceRequirements{},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      bread.Name + "-vol",
							MountPath: "/root",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: bread.Name + "-vol",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: bread.Spec.Task.Path,
						},
					},
				},
			},
		},
	}
	if err := r.Client.Create(ctx, &sshPod); err != nil {
		return err
	}
	return nil
}
