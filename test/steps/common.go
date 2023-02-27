package steps

import (
	"context"
	"fmt"

	"github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BrokerSteps struct {
	gomega.Gomega
	K8sClient        client.Client
	Cfg              *rest.Config
	TargetNamespace  *corev1.Namespace
	managedResources []client.Object
	By               func(string, ...func())
}

type Timing []interface{}

func (b *BrokerSteps) Manage(managedResource client.Object) {
	b.managedResources = append(b.managedResources, managedResource)
}

func (b *BrokerSteps) Cleanup() {
	b.By(fmt.Sprintf("cleaning up %d managed resources", len(b.managedResources)), func() {
		for _, r := range b.managedResources {
			b.By(fmt.Sprintf("cleaning up managed %s %s/%s", r.GetObjectKind().GroupVersionKind().Kind, r.GetNamespace(), r.GetName()))
			err := b.K8sClient.Delete(context.Background(), r)
			if err != nil && !errors.IsNotFound(err) {
				b.By(fmt.Sprintf("failed to clean up %s %s/%s", r.GetObjectKind().GroupVersionKind().Kind, r.GetNamespace(), r.GetName()))
			} else {
				b.By(fmt.Sprintf("successfully cleaned up %s %s/%s", r.GetObjectKind().GroupVersionKind().Kind, r.GetNamespace(), r.GetName()))
			}
		}

		b.managedResources = []client.Object{}
	})
}
