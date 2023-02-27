package test

import (
	ginkgo "github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Broker Lifecycle", func() {
	ginkgo.It("should run the full Broker Lifecycle", func() {
		brokerInstance := testSteps.GivenABroker()
		testSteps.ThenTheBrokerDeploymentIsReady(brokerInstance)

		testSteps.WhenTheBrokerResourceIsDeleted(brokerInstance)
		testSteps.ThenTheBrokerResourcesAreRemoved(brokerInstance)
	})
})
