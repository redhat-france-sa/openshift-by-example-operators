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

// CreateSecretForMongoDB initializes a new Secret for holding connection username and password.
func CreateSecretForMongoDB(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-mongodb-connectio",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "mongodb",
			},
		},
		StringData: map[string]string{
			"username":      spec.MongoDB.Username,
			"password":      spec.MongoDB.Password,
			"adminPassword": spec.MongoDB.Password,
		},
	}
}

// CreateDeploymentForMongoDB initializes a new Deployment for MongoDB.
func CreateDeploymentForMongoDB(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-mongodb",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "mongodb",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":              spec.AppName,
						"deploymentconfig": "mongodb",
						"container":        "mongodb",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "mongodb",
							Image:           spec.MongoDB.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "mongodb",
									ContainerPort: 27017,
									Protocol:      "TCP",
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-i",
											"-c",
											"mongo 127.0.0.1:27017/$MONGODB_DATABASE -u $MONGODB_USER -p $MONGODB_PASSWORD --eval=quit()",
										},
									},
								},
								InitialDelaySeconds: 3,
								FailureThreshold:    10,
								TimeoutSeconds:      1,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: int32(27017),
										},
									},
								},
								TimeoutSeconds:      1,
								InitialDelaySeconds: 30,
							},
							Env: []corev1.EnvVar{
								{
									Name: "MONGODB_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "username",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: spec.AppName + "-mongodb",
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
								{
									Name: "MONGODB_ADMIN_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											Key: "adminPassword",
											LocalObjectReference: corev1.LocalObjectReference{
												Name: spec.AppName + "-mongodb-connection",
											},
										},
									},
								},
								{
									Name:  "MONGODB_DATABASE",
									Value: spec.AppName,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      spec.AppName + "-mongodb-data",
									MountPath: "/var/lib/mongodb/data",
								},
							},
						},
					},
				},
			},
		},
	}

	if spec.MongoDB.Persistent {
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: spec.AppName + "-mongodb-data",
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: spec.AppName + "-mongodb",
					},
				},
			},
		}
	} else {
		deployment.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: spec.AppName + "-mongodb-data",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: corev1.StorageMediumDefault,
					},
				},
			},
		}
	}

	return deployment
}

// CreateServiceForMongoDB initializes a new Service for MongoDB.
func CreateServiceForMongoDB(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *corev1.Service {
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
				"container": "mongodb",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "mongodb",
					Port: 27017,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(27017),
					},
					Protocol: "TCP",
				},
			},
			Selector: map[string]string{
				"app":       spec.AppName,
				"container": "mongodb",
			},
			Type:            corev1.ServiceTypeClusterIP,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
