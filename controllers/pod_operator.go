package controllers

import (
	"context"
	corev1alpha1 "github.com/NJUPT-ISL/Breakfast/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

func GetPodLabel(bread *corev1alpha1.Bread) map[string]string {
	labels := map[string]string{"bread": bread.Name}
	if bread.Spec.Scv.Clock != "" {
		labels["scv/clock"] = bread.Spec.Scv.Clock
	}
	if bread.Spec.Scv.Priority != "" {
		labels["scv/priority"] = bread.Spec.Scv.Priority
	}
	if bread.Spec.Scv.Gpu != "0" {
		labels["scv/number"] = bread.Spec.Scv.Gpu
		if bread.Spec.Scv.Memory != "" {
			labels["scv/memory"] = bread.Spec.Scv.Memory
		}
	}

	return labels
}

func PodSchedulingSelector(bread *corev1alpha1.Bread) string {
	if bread.Spec.Scv.Gpu != "0" {
		return "yoda-scheduler2"
	}
	return "default-scheduler"
}

func TaskIsSSH(bread *corev1alpha1.Bread) bool {
	return bread.Spec.Task.Type == "ssh"
}

func GetPodImage(bread *corev1alpha1.Bread) string {
	if bread.Spec.Scv.Gpu != "0" {
		return "ccr.ccs.tencentyun.com/njupt-isl/" +
			bread.Spec.Framework.Name +
			"-gpu:" +
			bread.Spec.Framework.Version
	}
	return "ccr.ccs.tencentyun.com/njupt-isl/" +
		bread.Spec.Framework.Name +
		"-cpu:" +
		bread.Spec.Framework.Version
}

func (r *BreadReconciler) CreateSSHPodWithoutNodeSelected(ctx context.Context, bread *corev1alpha1.Bread) error {
	var sharePID = true
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
			ShareProcessNamespace: &sharePID,
			SchedulerName:         PodSchedulingSelector(bread),
			RestartPolicy:         v1.RestartPolicyNever,

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
					Ports: []v1.ContainerPort{
						{
							Name:          "ssh",
							ContainerPort: 22,
						},
					},
					Resources: v1.ResourceRequirements{},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      bread.Name + "-vol",
							MountPath: "/root",
						},
						{
							MountPath: "/dev/shm",
							Name: "dshm",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: bread.Name + "-vol",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/gluster-vol/" + bread.Namespace,
						},
					},
				},
				{
					Name: "dshm",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{
							Medium: "Memory",
						},
					},
				},
			},
		},
	}
	return r.Client.Create(ctx, &sshPod)
}

func (r *BreadReconciler) CreateSSHPodWithNodeSelected(ctx context.Context, bread *corev1alpha1.Bread, nodeName string) error {
	var sharePID = true
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
			ShareProcessNamespace: &sharePID,
			SchedulerName:         PodSchedulingSelector(bread),
			NodeName:              nodeName,
			RestartPolicy:         v1.RestartPolicyNever,
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
					Ports: []v1.ContainerPort{
						{
							Name:          "ssh",
							ContainerPort: 22,
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
							Path: "/gluster-vol/" + bread.Namespace,
						},
					},
				},
				{
					Name: "dshm",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{
							Medium: "Memory",
						},
					},
				},
			},
		},
	}
	return r.Client.Create(ctx, &sshPod)
}

func (r *BreadReconciler) CreateTaskPod(ctx context.Context, bread *corev1alpha1.Bread) error {
	var sharePID = true
	var TaskPod = v1.Pod{
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
			ShareProcessNamespace: &sharePID,
			SchedulerName:         PodSchedulingSelector(bread),
			RestartPolicy:         v1.RestartPolicyNever,
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
					Command:   strings.Split(bread.Spec.Task.Command, " "),
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
							Path: "/gluster-vol/" + bread.Namespace,
						},
					},
				},
			},
		},
	}
	return r.Client.Create(ctx, &TaskPod)
}

func (r *BreadReconciler) DeletePod(ctx context.Context, req ctrl.Request) error {
	pod := v1.Pod{}
	err := r.Client.Get(ctx, req.NamespacedName, &pod)
	if err != nil {
		return err
	}
	if err = r.Delete(ctx, &pod); err != nil {
		return err
	}
	return nil
}
