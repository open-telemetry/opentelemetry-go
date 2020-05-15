namespace-k8s:
	k apply -f k8s/namespace.yaml

jaeger-operator-k8s:
	# Create the jaeger operator and necessary artifacts in ns observability
	kubectl apply -f k8s/jaeger/jaegertracing.io_jaegers_crd.yaml
	kubectl apply -f k8s/jaeger/jaeger_service_account.yaml
	kubectl apply -f k8s/jaeger/jaeger_role.yaml
	kubectl apply -f k8s/jaeger/jaeger_role_binding.yaml
	kubectl apply -f k8s/jaeger/jaeger_operator.yaml

	# Create the cluster role & bindings
	kubectl apply -f k8s/jaeger/jaeger_cluster_role.yaml
	kubectl apply -f k8s/jaeger/jaeger_cluster_role_binding.yaml

jaeger-k8s:
	kubectl apply -f k8s/jaeger/jaeger.yaml

otel-collector-k8s:
	kubectl apply -f k8s/otel-collector.yaml

clean-k8s:
	- kubectl delete -f k8s/otel-collector.yaml

	- kubectl delete -f k8s/jaeger/jaeger.yaml

	- kubectl delete -f k8s/jaeger/jaeger_cluster_role.yaml
	- kubectl delete -f k8s/jaeger/jaeger_cluster_role_binding.yaml
	
	- kubectl delete -f k8s/jaeger/jaegertracing.io_jaegers_crd.yaml
	- kubectl delete -f k8s/jaeger/jaeger_service_account.yaml
	- kubectl delete -f k8s/jaeger/jaeger_role.yaml
	- kubectl delete -f k8s/jaeger/jaeger_role_binding.yaml
	- kubectl delete -f k8s/jaeger/jaeger_operator.yaml