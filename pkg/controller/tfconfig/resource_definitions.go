package tfconfig

import (
	configv1alpha1 "github.com/Svimba/tungstenfabric-config-operator/pkg/apis/config/v1alpha1"
	betav1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getConfigMapsObject(cr *configv1alpha1.TFConfig) []corev1.EnvFromSource {
	var list []corev1.EnvFromSource

	for _, cfm := range cr.Status.ConfigMapList {
		cfobj := &corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cfm,
				},
			},
		}

		list = append(list, *cfobj)
	}
	return list
}

func getContainerPortObject(cr *configv1alpha1.TFConfig) []corev1.ContainerPort {
	var list []corev1.ContainerPort
	for _, p := range cr.Spec.APISpec.Ports {
		pobj := &corev1.ContainerPort{
			Name:          p.Name,
			ContainerPort: p.Port,
		}
		list = append(list, *pobj)
	}
	return list
}

func getServicePortObject(cr *configv1alpha1.TFConfig) []corev1.ServicePort {
	var list []corev1.ServicePort
	for _, p := range cr.Spec.APISpec.Ports {
		pobj := &corev1.ServicePort{
			Name: p.Name,
			Port: p.Port,
		}
		list = append(list, *pobj)
	}
	return list
}

// newSTSForAPI retuns statefulset definition
func newDeploymentForAPI(cr *configv1alpha1.TFConfig) *betav1.Deployment {
	labels := map[string]string{
		"app": cr.Name + "-api",
	}

	deploy := &betav1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-api",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: betav1.DeploymentSpec{
			Replicas: cr.Spec.APISpec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cr.Name + "-api",
							Image: cr.Spec.APISpec.Image,
						},
					},
				},
			},
		},
	}
	// set co config maps
	configMaps := getConfigMapsObject(cr)
	if len(configMaps) > 0 {
		deploy.Spec.Template.Spec.Containers[0].EnvFrom = configMaps
	}
	// set ports if defined
	ports := getContainerPortObject(cr)
	if len(ports) > 0 {
		deploy.Spec.Template.Spec.Containers[0].Ports = ports
	} else {
		// If ports are not defined
		deploy.Spec.Template.Spec.Containers[0].Ports = []corev1.ContainerPort{
			{Name: "api", ContainerPort: 9100},
			{Name: "introspect", ContainerPort: 8084},
		}
	}

	return deploy
}

func newServicesForAPI(cr *configv1alpha1.TFConfig) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name + "-api",
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-api",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
		},
	}
	ports := getServicePortObject(cr)
	if len(ports) > 0 {
		service.Spec.Ports = ports
	} else {
		// If ports are not defined
		service.Spec.Ports = []corev1.ServicePort{
			{Name: "api", Port: 9100},
			{Name: "introspect", Port: 8084},
		}
	}
	return service
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
// func newPodForCR(cr *configv1alpha1.TFConfig) *corev1.Pod {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	return &corev1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      cr.Name + "-pod",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    "busybox",
// 					Image:   "busybox",
// 					Command: []string{"sleep", "3600"},
// 				},
// 			},
// 		},
// 	}
// }
