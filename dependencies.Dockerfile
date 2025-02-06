# This is a renovate-friendly source of Docker images.
FROM python:3.13.1-slim-bullseye@sha256:0eab754489555b1f0a6beaa0733e5bc7d39b12c58cb3f0a45613bec60555e716 AS python
FROM otel/weaver:v0.12.0@sha256:0b6136dc8ba68b3ee143dc8ee63b214af740276e5bbb0e7712ad61acc9b447da AS weaver
