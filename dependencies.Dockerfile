# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:81b94d27c19bba9f182fa3e46f13e21e01c48b8f5725972d82bab4cbe1bb96a2 AS python
FROM otel/weaver:v0.13.2@sha256:ae7346b992e477f629ea327e0979e8a416a97f7956ab1f7e95ac1f44edf1a893 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
