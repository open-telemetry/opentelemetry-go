# This is a renovate-friendly source of Docker images.
FROM python:3.13.5-slim-bullseye@sha256:77cabcb54d2633f258b3ed4c19d779918b385908262e27ef8d8067eb68fb527a AS python
FROM otel/weaver:v0.15.3@sha256:a84032d6eb95b81972d19de61f6ddc394a26976c1c1697cf9318bef4b4106976 AS weaver
FROM avtodev/markdown-lint:v1@sha256:6aeedc2f49138ce7a1cd0adffc1b1c0321b841dc2102408967d9301c031949ee AS markdown
