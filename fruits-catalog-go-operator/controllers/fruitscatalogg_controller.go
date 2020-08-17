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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	redhatcomv1alpha1 "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator/api/v1alpha1"
	deployment "github.com/redhat-france-sa/openshift-by-example-operators/fruits-catalog-go-operator/controllers/deployment"

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
// +kubebuilder:rbac:groups=core,resources=secrets;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile reads that state of the cluster for a FruitsCatalogG object and makes changes based on the state read
// and what is in the FruitsCatalogGSpec.
func (r *FruitsCatalogGReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("fruitscatalogg", req.NamespacedName)

	// your logic here

	// Fetch the  FruitsCatalogG instance of this reconcile request.
	instance := &redhatcomv1alpha1.FruitsCatalogG{}
	err := r.Get(ctx, req.NamespacedName, instance)
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
		mongodbSecret := deployment.CreateSecretForMongoDB(&instance.Spec, instance.Namespace)
		// Set App instance as the owner and controller.
		// NOTE: calling SetControllerReference, and setting owner references in
		// general, is important as it allows deleted objects to be garbage collected.
		controllerutil.SetControllerReference(instance, mongodbSecret, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbSecret); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Create " + mongodbSecret.Name + " Secret for MongoDB connection details")
		}
		mongodbDeployment := deployment.CreateDeploymentForMongoDB(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, mongodbDeployment, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbDeployment); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Apply " + mongodbDeployment.Name + " Deployment for MongoDB")
		}
		mongodbService := deployment.CreateServiceForMongoDB(&instance.Spec, instance.Namespace)
		controllerutil.SetControllerReference(instance, mongodbService, r.Scheme)
		if err := r.Client.Create(context.TODO(), mongodbService); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		} else if err == nil {
			reqLogger.Info("Apply " + mongodbService.Name + " Service for MongoDB")
		}
	}

	webappDeployment := deployment.CreateDeploymentForWebapp(&instance.Spec, instance.Namespace)
	controllerutil.SetControllerReference(instance, webappDeployment, r.Scheme)
	if err := r.Client.Create(context.TODO(), webappDeployment); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	} else if err == nil {
		reqLogger.Info("Apply " + webappDeployment.Name + " Deployment for WebApp")
	}
	webappService := deployment.CreateServiceForWebapp(&instance.Spec, instance.Namespace)
	controllerutil.SetControllerReference(instance, webappService, r.Scheme)
	if err := r.Client.Create(context.TODO(), webappService); err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	} else if err == nil {
		reqLogger.Info("Apply " + webappService.Name + " Service for WebApp")
	}

	// end of logic block

	return ctrl.Result{}, nil
}

// SetupWithManager takes care of controller initialization with Manager.
func (r *FruitsCatalogGReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redhatcomv1alpha1.FruitsCatalogG{}).
		Complete(r)
}

func (r *FruitsCatalogGReconciler) finalizeFruitsCatalogG(reqLogger logr.Logger, m *redhatcomv1alpha1.FruitsCatalogG) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
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
