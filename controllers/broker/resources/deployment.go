package resources

import (
	"fmt"

	mqttv1alpha1 "github.com/slivero/mqtt-broker-operator/api/v1alpha1"
	"github.com/slivero/mqtt-broker-operator/controllers/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func EnsureDeployment(ctx common.R12nContext, instance mqttv1alpha1.Broker) (err error) {

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	err = createOrUpdateDeployment(ctx, instance, &deployment, func() (err error) {
		deployment.Spec = appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"Broker": instance.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"Broker": instance.Name,
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-config", instance.Name),
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "mosquitto",
							Image: "eclipse-mosquitto:2.0.15",
							Command: []string{
								"mosquitto",
							},
							Args: []string{
								"-c",
								"/mosquitto/config/mosquitto.conf",
							},
							WorkingDir: "",
							Ports: []corev1.ContainerPort{
								{
									Name:          "mqtt",
									ContainerPort: 1883,
								},
								{
									Name:          "websockets",
									ContainerPort: 9001,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   map[corev1.ResourceName]resource.Quantity{"cpu": instance.Spec.Cores, "memory": instance.Spec.Memory},
								Requests: map[corev1.ResourceName]resource.Quantity{"cpu": instance.Spec.Cores, "memory": instance.Spec.Memory},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(1883),
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      5,
								PeriodSeconds:       5,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(1883),
									},
								},
								InitialDelaySeconds: 5,
								TimeoutSeconds:      5,
								PeriodSeconds:       5,
								SuccessThreshold:    3,
								FailureThreshold:    1,
							},
							ImagePullPolicy: "IfNotPresent",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/mosquitto/config/mosquitto.conf",
									SubPath:   "mosquitto.conf",
								},
							},
						},
					},
				},
			},
		}

		return
	})

	return
}

func createOrUpdateDeployment(ctx common.R12nContext, instance mqttv1alpha1.Broker, deployment *appsv1.Deployment, mutateFn func() (err error)) (err error) {
	result, err := ctrl.CreateOrUpdate(ctx.Context, ctx.Client, deployment, func() (err error) {
		err = mutateFn()
		if err != nil {
			return err
		}
		err = ctx.SetControllerReference(deployment)
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
		ctx.Log.Info("Deployment Created", "Name", deployment.Name, "Namespace", deployment.Namespace)

	case controllerutil.OperationResultUpdated:
		ctx.Log.Info("Deployment Created", "Name", deployment.Name, "Namespace", deployment.Namespace)
	default:
	}
	return
}
