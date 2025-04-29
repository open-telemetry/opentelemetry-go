# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:9fde509d8e79bdcd0d4aa735ac4a58ed6c9cc947d5a98b36eb1184144eeec2b1 AS python
FROM otel/weaver:v0.14.0@sha256:bea89bc5544ad760db2fd906c5285c2a3769c61fb04f660f9c31e7e44f11804b AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
