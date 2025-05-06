# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:d3f1e48b3e62e0e24b8ed20937d052662906c16e53013f32be88e2eb4f1b3532 AS python
FROM otel/weaver:v0.15.0@sha256:1cf1c72eaed57dad813c2e359133b8a15bd4facf305aae5b13bdca6d3eccff56 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
