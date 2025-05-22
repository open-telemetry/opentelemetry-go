# This is a renovate-friendly source of Docker images.
FROM python:3.13.3-slim-bullseye@sha256:45338d24f0fd55d4a7cb0cd3100a12a4350c7015cd1ec983b6f76d3d490a12c8 AS python
FROM otel/weaver:v0.15.0@sha256:1cf1c72eaed57dad813c2e359133b8a15bd4facf305aae5b13bdca6d3eccff56 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
