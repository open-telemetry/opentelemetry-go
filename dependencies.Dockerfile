# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:d3f1e48b3e62e0e24b8ed20937d052662906c16e53013f32be88e2eb4f1b3532 AS python
FROM otel/weaver:v0.14.0@sha256:bea89bc5544ad760db2fd906c5285c2a3769c61fb04f660f9c31e7e44f11804b AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
