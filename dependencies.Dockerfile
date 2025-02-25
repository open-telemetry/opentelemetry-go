# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:d3852c9e13043acf982e420e0e5f16b64a15223eae5e29d8526d232ba17e5cfc AS python
FROM otel/weaver:v0.13.2@sha256:ae7346b992e477f629ea327e0979e8a416a97f7956ab1f7e95ac1f44edf1a893 AS weaver
