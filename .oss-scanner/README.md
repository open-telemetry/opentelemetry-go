# OSS Scanner environment

This image follows the [OSS Scanner Dockerfile template][template]. It keeps the
checkout at `/src`, downloads dependencies for all Go modules (including
development tools), and compiles production code and tests. Go's default debug
information, module cache, and build cache remain available for reproductions,
fuzzing, and race detection after network access is removed.

Use these paths in the project's enrollment in
[anthropics/oss-scanner](https://github.com/anthropics/oss-scanner):

```yaml
dockerfile: .oss-scanner/Dockerfile
threat_model: THREAT_MODEL.md
```

Build from the repository root:

```sh
docker build --file .oss-scanner/Dockerfile --tag otel-go-oss-scanner .
```

Run fresh tests for all production modules, with the race detector and no network:

```sh
docker run --rm --network none \
  --env GOPROXY=off --env GOSUMDB=off \
  otel-go-oss-scanner \
  make SHELL=/bin/bash '.SHELLFLAGS=-e -o pipefail -c' \
  test-race ARGS='-race -count=1'
```

The Dockerfile runs a bounded short test suite during setup. Following the
scanner template, failures are logged without discarding the image; the GitHub
workflow requires the offline tests to pass. Tests can use loopback networking
inside the container, but cannot reach external services.

For interactive investigation, replace the `make` command with `bash`. The
shared OpenTelemetry security policy is saved at
`/opt/oss-scanner/SECURITY.md` during the build for offline reference. The
overrides in [THREAT_MODEL.md](../THREAT_MODEL.md) take precedence over that copy.

[template]: https://github.com/anthropics/oss-scanner/blob/main/templates/Dockerfile
