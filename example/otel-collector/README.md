# OpenTelemetry Collector Traces Example

This example illustrates how to export traces from the OpenTelemetry-Go SDK to the OpenTelemetry Collector, and from there to a Jaeger instance.
The complete flow is:

`Demo App -> OpenTelemetry Collector -> Jaeger`

# Prerequisites

The demo is built on Kubernetes, and uses a local instance of [microk8s](https://microk8s.io/). You will need access to a cluster in order to deploy the OpenTelemetry Collector and Jaeger components from this demo.

For simplicity, the demo application is not part of the k8s cluster, and will access the OpenTelemetry Collector through a NodePort on the cluster. Note that the NodePort opened by this demo is not secured. 

Ideally you'd want to either have your application running as part of the kubernetes cluster, or use a secured connection (NodePort/LoadBalancer with TLS or an ingress extension).

# Deploying Jaeger and OpenTelemetry Collector
The first step of this demo is to deploy a Jaeger instance and a Collector to your cluster. All the necessary Kubernetes deployment files are available in this demo, in the [k8s](./k8s) folder.
There are two ways to create the necessary deployments for this demo: using the [makefile](./Makefile) or manually applying the k8s files.

## Using the makefile

For using the [makefile](./Makefile), run the following commands in order:
```bash
# Create the namespace
make namespace-k8s

# Deploy Jaeger operator
make jaeger-operator-k8s

# After the operator is deployed, create the Jaeger instance
make jaeger-k8s

# Finally, deploy the OpenTelemetry Collector
make otel-collector-k8s
```

If you want to clean up after this, you can use the `make clean-k8s` to delete all the resources created above. Note that this will not remove the namespace. Because Kubernetes sometimes gets stuck when removing namespaces, please remove this namespace manually after all the resources inside have been deleted.

## Manual deployment

For manual deployments, follow the same steps as above, but instead run the `kubectl apply` yourself.

First, the namespace needs to be created:
```bash
k apply -f k8s/namespace.yaml
```

Jaeger is then deployed via the operator, and the demo follows [these steps](https://github.com/jaegertracing/jaeger-operator#getting-started) to create it:
```bash
# Create the jaeger operator and necessary artifacts in ns observability
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml

# Create the cluster role & bindings
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role.yaml
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role_binding.yaml

# Create the Jaeger instance itself:
kubectl apply -f k8s/jaeger/jaeger.yaml
```

The OpenTelemetry Collector is contained in a single k8s file, which can be deployed with one command: 
```bash 
k8s apply -f k8s/otel-collector.yaml
```


# Configuring the OpenTelemetry Collector

Although the above steps should deploy and configure both Jaeger and the OpenTelemetry Collector, it might be worth spending some time on the [configuration](./k8s/otel-collector.yaml) of the Collector.

One important part here is that, in order to enable our application to send traces to the OpenTelemetry Collector, we need to first configure the otlp receiver:

```yml
...
  otel-collector-config: |
    receivers:
      # Make sure to add the otlp receiver. 
      # This will open up the receiver on port 55680.
      otlp:
        endpoint: 0.0.0.0:55680
    processors:
...
```

This will create the receiver on the Collector side, and open up port `55680` for receiving traces.

The rest of the configuration is quite standard, with the only mention that we need to create the Jaeger exporter:

```yml
...
    exporters:
      jaeger_grpc:
        endpoint: "jaeger-collector.observability.svc.cluster.local:14250"
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
    port: 55680
    protocol: TCP
    targetPort: 55680
    nodePort: 30080
  - name: metrics # Default endpoint for metrics.
    port: 8888
    protocol: TCP
    targetPort: 8888
  selector:
    component: otel-collector
  type:
    NodePort
```

This service will bind the `55680` port used to access the otlp receiver to port `30080` on your cluster's node. By doing so, it makes it possible for us to access the Collector by using the static address `<node-ip>:30080`. In case you are running a local cluster, this will be `localhost:30080`. Note that you can also change this to a LoadBalancer or have an ingress extension for accessing the service.


# Writing the demo

Having the OpenTelemetry Collector started with the otlp port open for traces, and connected to Jaeger, let's look at the go app that will send traces to the Collector.

First, we need to create an exporter using the [otlp](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc) package:
```go
exp, err := otlp.NewExporter(otlp.WithInsecure(),
        // use the address of the NodePort service created above
        // <node-ip>:30080
        otlp.WithAddress("localhost:30080"), 
        otlp.WithGRPCDialOption(grpc.WithBlock()))
if err != nil {
        log.Fatalf("Failed to create the collector exporter: %v", err)
}
defer func() {
        err := exp.Stop()
        if err != nil {
                log.Fatalf("Failed to stop the exporter: %v", err)
        }
}()
```
This will initialize the exporter and connect to the otlp receiver at the address that we set for the [NodePort](#opentelemetry-collector-service): `localhost:30080`.

Feel free to remove the blocking operation, but it might come in handy when testing the connection.
Also, make sure to close the exporter before the app exits.

The next step is to create the TraceProvider:
```go
tp, err := sdktrace.NewProvider(
        sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
        sdktrace.WithResource(resource.New(
                // the service name used to display traces in Jaeger
                kv.Key(conventions.AttributeServiceName).String("test-service"),
        )),
        sdktrace.WithSyncer(exp))
if err != nil {
        log.Fatalf("error creating trace provider: %v\n", err)
}
```

It is important here to set the [AttributeServiceName](https://github.com/open-telemetry/opentelemetry-collector/blob/master/translator/conventions/opentelemetry.go#L20) from the `github.com/open-telemetry/opentelemetry-collector/translator/conventions` package on the resource level. This will be passed to the OpenTelemetry Collector, and used as ServiceName when exporting the traces to Jaeger.

After this, you can simply start sending traces:
```go
tracer := tp.Tracer("test-tracer")
ctx, span := tracer.Start(context.Background(), "CollectorExporter-Example")
defer span.End()
```

The traces should now be visible from the Jaeger UI (if you have it installed), or thorough the jaeger-query service, under the name `test-service`.

You can find the complete code for this example in the [main.go](./main.go) file.