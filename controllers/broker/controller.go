/*
Copyright 2023.

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

package broker

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	mqttv1alpha1 "github.com/slivero/mqtt-broker-operator/api/v1alpha1"
	"github.com/slivero/mqtt-broker-operator/controllers/broker/resources"
	"github.com/slivero/mqtt-broker-operator/controllers/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BrokerReconciler reconciles a Broker object
type BrokerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mosquitto.oliversmith.io,resources=brokers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mosquitto.oliversmith.io,resources=brokers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mosquitto.oliversmith.io,resources=brokers/finalizers,verbs=update
//+kubebuilder:rbac:groups=*,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments;configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Broker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *BrokerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	instance := &mqttv1alpha1.Broker{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// If the Broker isn't found then it must have been deleted so we don't need to do anything else
			// All child resources are cleaned up automatically
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	r12nContext := r.newContext(req, instance, ctx)

	err = resources.EnsureConfigMap(r12nContext, *instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = resources.EnsureDeployment(r12nContext, *instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BrokerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mqttv1alpha1.Broker{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *BrokerReconciler) newContext(req ctrl.Request, cr *mqttv1alpha1.Broker, ctx context.Context) common.R12nContext {
	return common.R12nContext{
		Context: ctx,
		Client:  r.Client,
		Log: r.Log.WithValues(
			"requestName", req.Name,
			"requestNamespace", req.Namespace,
		),
		SetControllerReference: func(controlled metav1.Object) error {
			return ctrl.SetControllerReference(cr, controlled, r.Scheme)
		},
	}
}
