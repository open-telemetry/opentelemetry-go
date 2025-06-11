# This is a renovate-friendly source of Docker images.
FROM python:3.13.4-slim-bullseye@sha256:473d423994c8a085132b3dab08c0848364d8b9c2240de1256f81c69581cfeb5b AS python
FROM otel/weaver:v0.15.2@sha256:b13acea09f721774daba36344861f689ac4bb8d6ecd94c4600b4d590c8fb34b9 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
