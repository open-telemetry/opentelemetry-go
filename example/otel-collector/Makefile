namespace-k8s:
	kubectl apply -f k8s/namespace.yaml

jaeger-operator-k8s:
	# Create the jaeger operator and necessary artifacts in ns observability
	kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/crds/jaegertracing.io_jaegers_crd.yaml
	kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/service_account.yaml
	kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/role.yaml
	kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/role_binding.yaml
	kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/operator.yaml

	# Create the cluster role & bindings
	kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/cluster_role.yaml
	kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/cluster_role_binding.yaml

jaeger-k8s:
	kubectl apply -f k8s/jaeger.yaml

prometheus-k8s:
	kubectl apply -f k8s/prometheus-service.yaml   # Prometheus instance
	kubectl apply -f k8s/prometheus-monitor.yaml   # Service monitor

otel-collector-k8s:
	kubectl apply -f k8s/otel-collector.yaml

clean-k8s:
	- kubectl delete -f k8s/otel-collector.yaml

	- kubectl delete -f k8s/prometheus-monitor.yaml
	- kubectl delete -f k8s/prometheus-service.yaml

	- kubectl delete -f k8s/jaeger.yaml

	- kubectl delete -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/cluster_role.yaml
	- kubectl delete -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/cluster_role_binding.yaml

	- kubectl delete -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/crds/jaegertracing.io_jaegers_crd.yaml
	- kubectl delete -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/service_account.yaml
	- kubectl delete -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/role.yaml
	- kubectl delete -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/role_binding.yaml
	- kubectl delete -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/v1.28.0/deploy/operator.yaml
