# This is a renovate-friendly source of Docker images.
FROM python:3.13.1-slim-bullseye AS python
FROM otel/weaver:v0.12.0 AS weaver
