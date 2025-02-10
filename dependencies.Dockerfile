# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:561ff65b26571534bea164cff88489f8ba621032475a099e572a9ccd4fbcd6ab AS python
FROM otel/weaver:v0.13.0@sha256:55fded24477d13e33ef501ffe192a78ec40c15962516712946b200f27451a493 AS weaver
