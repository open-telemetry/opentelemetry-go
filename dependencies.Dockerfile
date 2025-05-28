# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:f0acec66ba3578f142b7efc6a3e5ce62a5f51639ea430df576f9249c7baf6ef2 AS python
FROM otel/weaver:v0.15.1@sha256:95c0aaa493d84ac72a1188756bd46eec1ead8e82004e7778ff5779736be8d578 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
