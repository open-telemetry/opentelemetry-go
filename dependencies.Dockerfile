# This is a renovate-friendly source of Docker images.
FROM python:3.13.6-slim-bullseye@sha256:e98b521460ee75bca92175c16247bdf7275637a8faaeb2bcfa19d879ae5c4b9a AS python
FROM otel/weaver:v0.24.2@sha256:d1fb16d279f39810c340fbbf1cf9e5e995a3a9cefa531938e9012437e3bc00c1 AS weaver
FROM davidanson/markdownlint-cli2:v0.23.1@sha256:f382ea4fdc949883e79de678009437fb40c339323654c7b0dd4d5221cda8ed20 AS markdown
