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

package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	redhatcomv1alpha1 "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator-0.19.2/api/v1alpha1"
	deployment "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator-0.19.2/controllers/deployment"

	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// FruitsCatalogGReconciler reconciles a FruitsCatalogG object
type FruitsCatalogGReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=redhat.com,resources=fruitscataloggs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=redhat.com,resources=fruitscataloggs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=redhat.com,resources=fruitscataloggs/finalizers,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets;services;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups='',resources=deployments/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete

// Reconcile reads that state of the cluster for a FruitsCatalogG object and makes changes based on the state read
// and what is in the FruitsCatalogGSpec.
func (r *FruitsCatalogGReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	//ctx := context.Background()
	reqLogger := r.Log.WithValues("fruitscatalogg", req.NamespacedName)

	// your logic here
	reqLogger.Info("Starting reconcile loop for " + req.NamespacedName.Name)
	// Fetch the  FruitsCatalogG instance of this reconcile request.
	instance := &redhatcomv1alpha1.FruitsCatalogG{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue.
			reqLogger.Info("FruitsCatalogG resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get FruitsCatalogG.")
		return ctrl.Result{}, err
	}

	// Check if the FruitsCatalogG instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isFruitsCatalogGMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isFruitsCatalogGMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), "finalizer.fruitscatalogg.redhat.com") {
			// Run finalization logic for finalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeFruitsCatalogG(reqLogger, instance); err != nil {
				return ctrl.Result{}, err
			}

			// Remove finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(instance, "finalizer.fruitscatalogg.redhat.com")
			err := r.Update(context.TODO(), instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Ensure all required resources are up-to-date.
	installMongodb := instance.Spec.MongoDB.Install
	if installMongodb {
		// Deal with MongoDB connection credentials.
		mongodbSecret := deployment.CreateSecretForMongoDB(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, mongodbSecret, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbSecret); err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Error while creating "+mongodbSecret.Name+" Secret")
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Create " + mongodbSecret.Name + " Secret for MongoDB connection details")
			instance.Status.Secret = mongodbSecret.Name + " is holding connection details to MongoDB"
		}

		// Deal with MongoDB persistentn volume.
		if instance.Spec.MongoDB.Persistent {
			mongodbPVC := deployment.CreatePersistentVolumeClaimMongoDB(&instance.Spec, instance.Namespace)
			controllerutil.SetControllerReference(instance, mongodbPVC, r.Scheme)
			if err := r.Client.Create(context.TODO(), mongodbPVC); err != nil && !errors.IsAlreadyExists(err) {
				reqLogger.Error(err, "Error while creating "+mongodbPVC.Name+" PVC")
				return ctrl.Result{}, err
			} else if err == nil {
				reqLogger.Info("Apply " + mongodbPVC.Name + " PVC for MongoDB")
			}
		}

		// Deal with MongoDB deployment and service.
		mongodbDeployment := deployment.CreateDeploymentForMongoDB(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, mongodbDeployment, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbDeployment); err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Error while creating "+mongodbDeployment.Name+" Deployment")
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Apply " + mongodbDeployment.Name + " Deployment for MongoDB")
			instance.Status.MongoDB = mongodbDeployment.Name + " is the Deployment for MongoDB"
		}
		mongodbService := deployment.CreateServiceForMongoDB(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, mongodbService, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbService); err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Error while creating "+mongodbService.Name+" Service")
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Apply " + mongodbService.Name + " Service for MongoDB")
		}
	}

	// Deal with WebApp deployment and service.
	webappDeployment := deployment.CreateDeploymentForWebapp(&instance.Spec, instance.Namespace)
	controllerutil.SetControllerReference(instance, webappDeployment, r.Scheme)
	if err := r.Client.Create(context.TODO(), webappDeployment); err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Error(err, "Error while creating "+webappDeployment.Name+" Deployment")
		return ctrl.Result{}, err
	} else if err != nil && errors.IsAlreadyExists(err) {
		// Check if we got the correct number of replicas.
		if *webappDeployment.Spec.Replicas != instance.Spec.WebApp.ReplicaCount {
			webappDeployment.Spec.Replicas = &instance.Spec.WebApp.ReplicaCount
			if err2 := r.Client.Update(context.TODO(), webappDeployment); err2 != nil {
				reqLogger.Error(err2, "Error while updating replicas in "+webappDeployment.Name+" Deployment")
				return ctrl.Result{}, err2
			} else if err2 == nil {
				reqLogger.Info("Update replicas " + webappDeployment.Name + " Deployment for WebApp")
				instance.Status.WebApp = webappDeployment.Name + " is the Deployment for WebApp"
			}
		}
	} else if err == nil {
		reqLogger.Info("Apply " + webappDeployment.Name + " Deployment for WebApp")
		instance.Status.WebApp = webappDeployment.Name + " is the Deployment for WebApp"
	}
	webappService := deployment.CreateServiceForWebapp(&instance.Spec, instance.Namespace)
	controllerutil.SetControllerReference(instance, webappService, r.Scheme)
	if err := r.Client.Create(context.TODO(), webappService); err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Error(err, "Error while creating "+webappService.Name+" Service")
		return ctrl.Result{}, err
	} else if err == nil {
		reqLogger.Info("Apply " + webappService.Name + " Service for WebApp")
	}

	// Finally check if we're on OpenShift of Vanilla Kubernetes.
	isOpenShift := false

	// The discovery package is used to discover APIs supported by a Kubernetes API server.
	config, err := ctrl.GetConfig()
	if err == nil && config != nil {
		dclient, err := getDiscoveryClient(config)
		if err == nil && dclient != nil {
			apiGroupList, err := dclient.ServerGroups()
			if err != nil {
				reqLogger.Info("Error while querying ServerGroups, assuming we're on Vanilla Kubernetes")
			} else {
				for i := 0; i < len(apiGroupList.Groups); i++ {
					if strings.HasSuffix(apiGroupList.Groups[i].Name, ".openshift.io") {
						isOpenShift = true
						reqLogger.Info("We detected being on OpenShift! Wouhou!")
						break
					}
				}
			}
		} else {
			reqLogger.Info("Cannot retrieve a DiscoveryClient, assuming we're on Vanilla Kubernetes")
		}
	}
	// Create a Route if we're on OpenShift ;-)
	if isOpenShift {
		appRoute := deployment.CreateRouteForWebapp(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, appRoute, r.Scheme)
		if err := r.Client.Create(context.TODO(), appRoute); err != nil && !errors.IsAlreadyExists(err) {
			reqLogger.Error(err, "Error while creating "+appRoute.Name+" Route")
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Apply " + appRoute.Name + " Route for WebApp")
		}

		// Maybe on next reconciliation loop ?
		reqLogger.Info("Looking for route " + appRoute.ObjectMeta.Name + " on namespace " + instance.Namespace)
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: appRoute.ObjectMeta.Name, Namespace: instance.Namespace}, appRoute)
		if err == nil {
			instance.Status.Route = appRoute.Status.Ingress[0].Host
		} else {
			reqLogger.Error(err, "Error while reading Route for getting its Status.Ingress[0].Host field")
		}
	}

	// Updating the Status that is modeled as a subresource. This way we can update
	// the status of our resources without increasing the ResourceGeneration metadata field.
	r.Status().Update(context.Background(), instance)

	// end of logic block

	return ctrl.Result{}, nil
}

// SetupWithManager takes care of controller initialization with Manager.
func (r *FruitsCatalogGReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redhatcomv1alpha1.FruitsCatalogG{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *FruitsCatalogGReconciler) finalizeFruitsCatalogG(reqLogger logr.Logger, m *redhatcomv1alpha1.FruitsCatalogG) error {
	// TODO(user): Add the cleanup steps that the operator needs to do before the CR
	// can be deleted. Examples of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	reqLogger.Info("Successfully finalized FruitsCatalogG")
	return nil
}

func (r *FruitsCatalogGReconciler) addFinalizer(reqLogger logr.Logger, m *redhatcomv1alpha1.FruitsCatalogG) error {
	reqLogger.Info("Adding Finalizer for the FruitsCatalogG")
	controllerutil.AddFinalizer(m, "finalizer.fruitscatalogg.redhat.com")

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update FruitsCatalogG with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// getDiscoveryClient returns a discovery client for the current reconciler
func getDiscoveryClient(config *rest.Config) (*discovery.DiscoveryClient, error) {
	return discovery.NewDiscoveryClientForConfig(config)
}
