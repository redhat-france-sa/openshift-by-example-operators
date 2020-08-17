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

package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	redhatcomv1alpha1 "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator/api/v1alpha1"
)

// CreateDeploymentForWebapp initializes a new Deployment for Webapp.
func CreateDeploymentForWebapp(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-webapp",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "webapp",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.WebApp.ReplicaCount,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app":              spec.AppName,
				"deploymentconfig": "webapp",
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":              spec.AppName,
						"deploymentconfig": "webapp",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "webapp",
							Image:           spec.WebApp.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"/work/application",
							},
							Args: []string{
								"-Dquarkus.http.host=0.0.0.0",
								"-Dquarkus.mongodb.connection-string=mongodb://$(MONGODB_USER):$(MONGODB_PASSWORD)@" + spec.AppName + "-mongodb:27017",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health/ready",
										Port: intstr.IntOrString{
											Type:   intstr.String,
											StrVal: "http",
										},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health/live",
										Port: intstr.IntOrString{
											Type:   intstr.String,
											StrVal: "http",
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "MONGODB_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "username",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: spec.AppName + "-mongodb-connection",
											},
										},
									},
								},
								{
									Name: "MONGODB_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "password",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: spec.AppName + "-mongodb-connection",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

// CreateServiceForWebapp initializes a new Service for MongoDB.
func CreateServiceForWebapp(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-mongodb",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "webapp",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "webapp",
					Port: 80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "http",
					},
					Protocol: "TCP",
				},
			},
			Selector: map[string]string{
				"app":              spec.AppName,
				"deploymentconfig": "webapp",
			},
			Type:            corev1.ServiceTypeClusterIP,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
