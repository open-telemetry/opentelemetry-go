# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:561ff65b26571534bea164cff88489f8ba621032475a099e572a9ccd4fbcd6ab AS python
FROM otel/weaver:v0.13.1@sha256:1e462a7c4070d71a711b3326f64a31837a8ee6321f21b86e724020da94cb7105 AS weaver
