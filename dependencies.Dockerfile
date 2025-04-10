# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:0d46ec7010093c2a30ae712c3d6fc9d3938ae8d31dcf38c14deee3e43f88e6ca AS python
FROM otel/weaver:v0.13.2@sha256:ae7346b992e477f629ea327e0979e8a416a97f7956ab1f7e95ac1f44edf1a893 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
