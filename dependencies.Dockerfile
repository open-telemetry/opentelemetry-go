# This is a renovate-friendly source of Docker images.
FROM python:3.13.2-slim-bullseye@sha256:561ff65b26571534bea164cff88489f8ba621032475a099e572a9ccd4fbcd6ab AS python
FROM otel/weaver:v0.12.0@sha256:0b6136dc8ba68b3ee143dc8ee63b214af740276e5bbb0e7712ad61acc9b447da AS weaver
