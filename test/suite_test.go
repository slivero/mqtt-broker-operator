package test

import (
	"context"
	"fmt"
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	mqttv1alpha1 "github.com/slivero/mqtt-broker-operator/api/v1alpha1"
	"github.com/slivero/mqtt-broker-operator/test/steps"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var g gomega.Gomega

func TestE2e(t *testing.T) {
	g = gomega.NewGomega(ginkgo.Fail)
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E2e Tests")
}

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var targetNamespace corev1.Namespace
var testSteps steps.BrokerSteps

var _ = ginkgo.BeforeSuite(func() {
	useExistingCluster := true
	testEnv = &envtest.Environment{
		UseExistingCluster: &useExistingCluster,
	}

	cfg, err := testEnv.Start()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(cfg).NotTo(gomega.BeNil())

	gomega.Expect(mqttv1alpha1.AddToScheme(scheme.Scheme)).To(gomega.Succeed())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(k8sClient).NotTo(gomega.BeNil())

	gomega.Expect(k8sClient.Get(context.Background(), types.NamespacedName{Name: "mqtt-broker-operator-system"}, &targetNamespace)).To(gomega.Succeed())

	ginkgo.By(fmt.Sprintf("Using Namespace %s", targetNamespace.Name))

})

var _ = ginkgo.AfterEach(func() {
	testSteps.Cleanup()
	testEnv.Stop()
})

var _ = ginkgo.BeforeEach(func() {
	testSteps = steps.BrokerSteps{
		Gomega:          g,
		K8sClient:       k8sClient,
		Cfg:             cfg,
		TargetNamespace: &targetNamespace,
		By:              ginkgo.By,
	}

	ginkgo.By(fmt.Sprintf("Using Namespace %s", targetNamespace.Name))
})
