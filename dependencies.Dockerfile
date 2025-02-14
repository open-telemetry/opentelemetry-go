# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:561ff65b26571534bea164cff88489f8ba621032475a099e572a9ccd4fbcd6ab AS python
FROM otel/weaver:v0.13.2@sha256:ae7346b992e477f629ea327e0979e8a416a97f7956ab1f7e95ac1f44edf1a893 AS weaver
