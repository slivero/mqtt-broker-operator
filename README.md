# mqtt-broker-operator

Simple Kubernetes operator to deploy / manage [Eclipse Mosquitto](https://mosquitto.org/), an MQTT message broker. 

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [kind](https://sigs.k8s.io/kind) to get a local cluster
 for testing, the makefile includes helper functions for this:

### Create a local kind cluster 

```sh
make deploy-local-cluster
```

### Running the operator locally
Build and install the operator into the cluster (uses the currently active kubectl context):

```sh
make local-build-and-deploy
```

### Deploying a Broker

Save and apply the yaml and apply it:

*broker.yaml*
```yaml
apiVersion: mosquitto.oliversmith.io/v1alpha1
kind: Broker
metadata:
  name: broker-sample
spec:
  cores: 1
  memory: 1Gi

```

```sh
kubectl apply -f broker.yaml
```

### Running Tests

A simple end to end test is included, this verifies a Mosquitto broker can be provisioned and then tears it down:

```sh
make test-e2e
```

### Cleaning Up

To tear down the local kind cluster:

```sh
make cleanup-local-cluster
```
