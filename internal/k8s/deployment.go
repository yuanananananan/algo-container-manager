package k8s

import (
	"algo-container-manager/internal/model"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateDeployment(clientset *kubernetes.Clientset, req model.StartAlgorithmRequest) error {
	deploymentName := req.DeploymentName
	if req.DeploymentName == "" {
		return fmt.Errorf("deployment name is required")
	}
	labels := map[string]string{
		"app":                  deploymentName,
		"algorithm-version-id": fmt.Sprintf("%d", req.AlgorithmVersionID),
		"image-id":             fmt.Sprintf("%d", req.ImageID),
	}

	replicas := req.Replicas
	if replicas == 0 {
		replicas = 1
	}

	envVars := []corev1.EnvVar{}
	for k, v := range req.Env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deploymentName,
							Image: req.Image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: req.Port},
							},
							Env: envVars,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(req.CPU),
									corev1.ResourceMemory: resource.MustParse(req.Memory),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(req.CPU),
									corev1.ResourceMemory: resource.MustParse(req.Memory),
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt32(req.Port),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt32(req.Port),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}
	_, err := clientset.AppsV1().
		Deployments(req.Namespace).
		Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("create deployment failed: %w", err)
	}
	return nil
}
