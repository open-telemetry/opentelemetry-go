#!/bin/bash

docker run \
-v "${PWD}/logging.yaml":/logging.yaml \
-p 4318:4318 \
otel/opentelemetry-collector \
--config logging.yaml;
