package steps

import (
	"context"
	"fmt"
	"strings"

	"github.com/onsi/gomega"
	mqttv1alpha1 "github.com/slivero/mqtt-broker-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/icrowley/fake"
)

func (s *BrokerSteps) GivenABroker() mqttv1alpha1.Broker {
	broker := mqttv1alpha1.Broker{
		ObjectMeta: metav1.ObjectMeta{
			Name:      aName(),
			Namespace: s.TargetNamespace.Name,
		},
		Spec: mqttv1alpha1.BrokerSpec{
			Memory: resource.MustParse("1Gi"),
			Cores:  resource.MustParse("100m"),
		},
	}

	s.By("creating Broker resource", func() {
		s.Expect(s.K8sClient.Create(context.Background(), &broker)).To(gomega.Succeed())
		s.Manage(&broker)
	})

	return broker
}

func (s *BrokerSteps) ThenTheBrokerDeploymentIsReady(brokerInstance mqttv1alpha1.Broker) {
	brokerDeployment := &appsv1.Deployment{}

	s.By("checking the broker deployment becomes ready", func() {
		s.Eventually(func() (replicasReady int32) {
			s.K8sClient.Get(context.Background(), types.NamespacedName{Name: brokerInstance.Name, Namespace: brokerInstance.Namespace}, brokerDeployment)
			return brokerDeployment.Status.ReadyReplicas
		}, 60, 1).Should(gomega.Equal(int32(1)))
	})
}

func (s *BrokerSteps) WhenTheBrokerResourceIsDeleted(brokerInstance mqttv1alpha1.Broker) {
	s.By("deleting Broker resource", func() {
		s.Expect(s.K8sClient.Delete(context.Background(), &brokerInstance)).To(gomega.Succeed())
	})
}

func (s *BrokerSteps) ThenTheBrokerResourcesAreRemoved(brokerInstance mqttv1alpha1.Broker) {

	s.By("checking the broker resource is removed", func() {
		s.Eventually(func() bool {
			broker := &mqttv1alpha1.Broker{}
			return errors.IsNotFound(s.K8sClient.Get(context.Background(), types.NamespacedName{Name: brokerInstance.Name, Namespace: brokerInstance.Namespace}, broker))
		}, 60, 1).Should(gomega.BeTrue())
	})

	s.By("checking the configmap resource is removed", func() {
		s.Eventually(func() bool {
			configmap := &mqttv1alpha1.Broker{}
			return errors.IsNotFound(s.K8sClient.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf("%s-config", brokerInstance.Name), Namespace: brokerInstance.Namespace}, configmap))
		}, 60, 1).Should(gomega.BeTrue())
	})
}

func aName() string {
	return strings.ToLower(fmt.Sprintf("%s-%s", fake.FirstName(), fake.LastName()))
}
