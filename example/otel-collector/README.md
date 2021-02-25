# OpenTelemetry Collector Traces Example

This example illustrates how to export trace and metric data from the
OpenTelemetry-Go SDK to the OpenTelemetry Collector. From there, we bring the
trace data to Jaeger and the metric data to Prometheus
The complete flow is:

```
                                          -----> Jaeger (trace)
App + SDK ---> OpenTelemetry Collector ---|
                                          -----> Prometheus (metrics)
```

# Prerequisites
You will need access to a Kubernetes cluster for this demo. We use a local
instance of [microk8s](https://microk8s.io/), but please feel free to pick
your favorite. If you do decide to use microk8s, please ensure that dns
and storage addons are enabled

```bash
microk8s enable dns storage
```

For simplicity, the demo application is not part of the k8s cluster, and will
access the OpenTelemetry Collector through a NodePort on the cluster. Note that
the NodePort opened by this demo is not secured.

Ideally you'd want to either have your application running as part of the
kubernetes cluster, or use a secured connection (NodePort/LoadBalancer with TLS
or an ingress extension).

# Deploying to Kubernetes
All the necessary Kubernetes deployment files are available in this demo, in the
[k8s](./k8s) folder. For your convenience, we assembled a [makefile](./Makefile)
with deployment commands (see below). For those with subtly different systems,
you are, of course, welcome to poke inside the Makefile and run the commands
manually. If you use microk8s and alias `microk8s kubectl` to `kubectl`, the
Makefile will not recognize the alias, and so the commands will have to be run
manually.

## Setting up the Prometheus operator
If you're using microk8s like us, simply do
```bash
microk8s enable prometheus
```
and you're good to go. Move on to [Using the makefile](#using-the-makefile).

Otherwise, obtain a copy of the Prometheus Operator stack from
[coreos](https://github.com/coreos/kube-prometheus):
```bash
git clone https://github.com/coreos/kube-prometheus.git
cd kube-prometheus
kubectl create -f manifests/setup

# wait for namespaces and CRDs to become available, then
kubectl create -f manifests/
```

And to tear down the stack when you're finished:
```bash
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
```

## Using the makefile
Next, we can deploy our Jaeger instance, Prometheus monitor, and Collector
using the [makefile](./Makefile).

```bash
# Create the namespace
make namespace-k8s

# Deploy Jaeger operator
make jaeger-operator-k8s

# After the operator is deployed, create the Jaeger instance
make jaeger-k8s

# Then the Prometheus instance. Ensure you have enabled a Prometheus operator
# before executing (see above).
make prometheus-k8s

# Finally, deploy the OpenTelemetry Collector
make otel-collector-k8s
```

If you want to clean up after this, you can use the `make clean-k8s` to delete
all the resources created above. Note that this will not remove the namespace.
Because Kubernetes sometimes gets stuck when removing namespaces, please remove
this namespace manually after all the resources inside have been deleted,
for example with

```bash
kubectl delete namespaces observability
```

# Configuring the OpenTelemetry Collector
Although the above steps should deploy and configure everything, let's spend
some time on the [configuration](./k8s/otel-collector.yaml) of the Collector.

One important part here is that, in order to enable our application to send data
to the OpenTelemetry Collector, we need to first configure the `otlp` receiver:

```yml
...
  otel-collector-config: |
    receivers:
      # Make sure to add the otlp receiver.
      # This will open up the receiver on port 4317.
      otlp:
        endpoint: 0.0.0.0:4317
    processors:
...
```

This will create the receiver on the Collector side, and open up port `4317`
for receiving traces.

The rest of the configuration is quite standard, with the only mention that we
need to create the Jaeger and Prometheus exporters:

```yml
...
    exporters:
      jaeger_grpc:
        endpoint: "jaeger-collector.observability.svc.cluster.local:14250"

      prometheus:
           endpoint: 0.0.0.0:8889
           namespace: "testapp"
...
```

## OpenTelemetry Collector service

One more aspect in the OpenTelemetry Collector [configuration](./k8s/otel-collector.yaml) worth looking at is the NodePort service used for accessing it:
```yaml
apiVersion: v1
kind: Service
metadata:
        ...
spec:
  ports:
  - name: otlp # Default endpoint for otlp receiver.
    port: 4317
    protocol: TCP
    targetPort: 4317
    nodePort: 30080
  - name: metrics # Endpoint for metrics from our app.
    port: 8889
    protocol: TCP
    targetPort: 8889
  selector:
    component: otel-collector
  type:
    NodePort
```

This service will bind the `55680` port used to access the otlp receiver to port `30080` on your cluster's node. By doing so, it makes it possible for us to access the Collector by using the static address `<node-ip>:30080`. In case you are running a local cluster, this will be `localhost:30080`. Note that you can also change this to a LoadBalancer or have an ingress extension for accessing the service.


# Running the code
You can find the complete code for this example in the [main.go](./main.go)
file. To run it, ensure you have a somewhat recent version of Go (preferably >=
1.13) and do

```bash
go run main.go
```

The example simulates an application, hard at work, computing for ten seconds
then finishing.

# Viewing instrumentation data
Now the exciting part! Let's check out the telemetry data generated by our
sample application

## Jaeger UI
First, we need to enable an ingress provider. If you've been using microk8s,
do

```bash
microk8s enable ingress
```

Then find out where the Jaeger console is living:
```bash
kubectl get ingress --all-namespaces
```

For us, we get the output
```
NAMESPACE       NAME           CLASS    HOSTS   ADDRESS     PORTS   AGE
observability   jaeger-query   <none>   *       127.0.0.1   80      5h40m
```
indicating that the Jaeger UI is available at
[http://localhost:80](http://localhost:80). Navigate there in your favorite
web-browser to view the generated traces.

## Prometheus
Unfortunately, the Prometheus operator doesn't provide a convenient
out-of-the-box ingress route for us to use, so we'll use port-forwarding
instead. Note: this is a quick-and-dirty solution for the sake of example.
You *will* be attacked by shady people if you do this in production!

```bash
kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090
```

Then navigate to [http://localhost:9090](http://localhost:9090) to view
the Prometheus dashboard.
