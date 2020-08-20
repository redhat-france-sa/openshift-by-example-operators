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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	redhatcomv1alpha1 "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator/api/v1alpha1"

	routev1 "github.com/openshift/api/route/v1"
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
								"-Dquarkus.mongodb.connection-string=mongodb://$(MONGODB_USER):$(MONGODB_PASSWORD)@" + spec.AppName + "-mongodb:27017/" + spec.AppName,
								"-Dquarkus.mongodb.database=" + spec.AppName,
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

// CreateServiceForWebapp initializes a new Service for WebApp.
func CreateServiceForWebapp(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName + "-webapp",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "webapp",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "http",
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

// CreateRouteForWebapp initializes a new Route for WebApp.
func CreateRouteForWebapp(spec *redhatcomv1alpha1.FruitsCatalogGSpec, namespace string) *routev1.Route {
	weight := int32(100)
	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.AppName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":       spec.AppName,
				"container": "webapp",
			},
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Name:   spec.AppName + "-webapp",
				Kind:   "Service",
				Weight: &weight,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "http",
				},
			},
			TLS: &routev1.TLSConfig{
				Termination:                   routev1.TLSTerminationEdge,
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
			},
			WildcardPolicy: routev1.WildcardPolicyNone,
		},
	}
}
