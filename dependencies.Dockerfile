# This is a renovate-friendly source of Docker images.
FROM python:3.13.5-slim-bullseye@sha256:ba65ee6bad4e448a9d7214bd3bef36ef908c05df601264c5e067816e18971ff6 AS python
FROM otel/weaver:v0.16.1@sha256:5ca4901b460217604ddb83feaca05238e2b016a226ecfb9b87a95555918a03af AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
