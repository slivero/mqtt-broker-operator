package resources

import (
	"fmt"

	mqttv1alpha1 "github.com/slivero/mqtt-broker-operator/api/v1alpha1"
	"github.com/slivero/mqtt-broker-operator/controllers/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func EnsureConfigMap(ctx common.R12nContext, instance mqttv1alpha1.Broker) (err error) {

	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-config", instance.Name),
			Namespace: instance.Namespace,
		},
	}

	err = createOrUpdateConfigMap(ctx, instance, &configMap, func() (err error) {
		configMap.Data = map[string]string{
			"mosquitto.conf": `
listener 1883
allow_anonymous true`,
		}
		return
	})

	return
}

func createOrUpdateConfigMap(ctx common.R12nContext, instance mqttv1alpha1.Broker, configMap *corev1.ConfigMap, mutateFn func() (err error)) (err error) {
	result, err := ctrl.CreateOrUpdate(ctx.Context, ctx.Client, configMap, func() (err error) {
		err = mutateFn()
		if err != nil {
			return err
		}
		err = ctx.SetControllerReference(configMap)
		if err != nil {
			return err
		}
		return
	})

	if err != nil {
		return
	}

	switch result {
	case controllerutil.OperationResultCreated:
		ctx.Log.Info("ConfigMap Created", "Name", configMap.Name, "Namespace", configMap.Namespace)

	case controllerutil.OperationResultUpdated:
		ctx.Log.Info("ConfigMap Created", "Name", configMap.Name, "Namespace", configMap.Namespace)
	default:
	}
	return

}
