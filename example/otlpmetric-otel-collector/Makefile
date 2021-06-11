namespace-k8s:
	kubectl apply -f k8s/namespace.yaml

prometheus-k8s:
	kubectl apply -f k8s/prometheus-service.yaml   # Prometheus instance
	kubectl apply -f k8s/prometheus-monitor.yaml   # Service monitor

otel-collector-k8s:
	kubectl apply -f k8s/otel-collector.yaml

clean-k8s:
	- kubectl delete -f k8s/otel-collector.yaml

	- kubectl delete -f k8s/prometheus-monitor.yaml
	- kubectl delete -f k8s/prometheus-service.yaml
