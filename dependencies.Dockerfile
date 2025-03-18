# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:c527a33e5265d0f830994d1b3237d38840a7b7986a8b9374a4b941ac34048190 AS python
FROM otel/weaver:v0.13.2@sha256:ae7346b992e477f629ea327e0979e8a416a97f7956ab1f7e95ac1f44edf1a893 AS weaver
