/*
MIT License

Copyright (c) 2020 RH France Solution Architects

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	redhatcomv1beta1 "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator-1.0.0/api/v1beta1"
)

// CreateSecretForMongoDB initializes a new Secret for holding connection username and password.
func CreateSecretForMongoDB(spec *redhatcomv1beta1.FruitsCatalogG1Spec, namespace string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-mongodb-connection",
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

// CreatePersistentVolumeClaimMongoDB initializes a new PerssitentVolumeClaim for MongoDB.
func CreatePersistentVolumeClaimMongoDB(spec *redhatcomv1beta1.FruitsCatalogG1Spec, namespace string) *corev1.PersistentVolumeClaim {
	claim := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(spec.MongoDB.VolumeSize),
				},
			},
		},
	}

	return claim
}

// CreateDeploymentForMongoDB initializes a new Deployment for MongoDB.
func CreateDeploymentForMongoDB(spec *redhatcomv1beta1.FruitsCatalogG1Spec, namespace string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-mongodb",
			Namespace: namespace,
			Labels: map[string]string{
				"app":              spec.AppName,
				"deploymentconfig": "mongodb",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app":              spec.AppName,
				"deploymentconfig": "mongodb",
				"container":        "mongodb",
			}},
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
											"mongo 127.0.0.1:27017/$MONGODB_DATABASE -u $MONGODB_USER -p $MONGODB_PASSWORD --eval=\"quit()\"",
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
func CreateServiceForMongoDB(spec *redhatcomv1beta1.FruitsCatalogG1Spec, namespace string) *corev1.Service {
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
