package tfconfig

import (
	"context"

	betav1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// API deployment handler
// create/update(if exists) config API deplyment
// return true/false(Requeue), error
func (r *ReconcileTFConfig) handleAPIDeployment() (bool, error) {
	// Define a new API deployment object
	apiDeployment := newDeploymentForAPI(r.instance)
	// Set TFConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(r.instance, apiDeployment, r.scheme); err != nil {
		return false, err
	}
	// Check if this API deployment already exists
	foundAPIDeploy := &betav1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: apiDeployment.Name, Namespace: apiDeployment.Namespace}, foundAPIDeploy)
	if err != nil && errors.IsNotFound(err) {
		r.reqLogger.Info("Creating a new API deployment", "Deploy.Namespace", apiDeployment.Namespace, "Deploy.Name", apiDeployment.Name)
		err = r.client.Create(context.TODO(), apiDeployment)
		if err != nil {
			return false, err
		}
		// Deployment has been created successfully - don't requeue
		return false, nil
	} else if err != nil {
		return false, err
	}
	// Deployment already exists - don't requeue
	r.reqLogger.Info("Skip reconcile: API Deployment already exists", "Deploy.Namespace", foundAPIDeploy.Namespace, "Deploy.Name", foundAPIDeploy.Name)
	return false, nil
}

// API service handler
// create/update(if exists) config API deplyment
// return true/false(Requeue), error
func (r *ReconcileTFConfig) handleAPIService() (bool, error) {
	// Define a new API service object
	apiService := newServicesForAPI(r.instance)
	// Set TFConfig instance as the owner and controller
	if err := controllerutil.SetControllerReference(r.instance, apiService, r.scheme); err != nil {
		return false, err
	}
	// Check if this API Service already exists
	foundAPIService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: apiService.Name, Namespace: apiService.Namespace}, foundAPIService)
	if err != nil && errors.IsNotFound(err) {
		r.reqLogger.Info("Creating a new API Service", "Service.Namespace", apiService.Namespace, "Service.Name", apiService.Name)
		err = r.client.Create(context.TODO(), apiService)
		if err != nil {
			return false, err
		}
		// Service has been created successfully - don't requeue
		return false, nil
	} else if err != nil {
		return false, err
	}
	// Service already exists - don't requeue
	r.reqLogger.Info("Skip reconcile: API Service already exists", "Service.Namespace", foundAPIService.Namespace, "Service.Name", foundAPIService.Name)
	return false, nil
}
