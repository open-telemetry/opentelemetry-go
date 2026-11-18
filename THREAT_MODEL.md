# Threat model

See the
[OpenTelemetry Security Policy](https://github.com/open-telemetry/.github/blob/main/SECURITY.md)
for the security model, severity assessment, and private vulnerability reporting
process. It applies to all artifacts in this repository, except for the explicit
override below.

## What this project does and where untrusted input enters

- OpenTelemetry Go provides APIs, SDKs, and exporters for traces, metrics, and
  logs in Go applications.
- Untrusted input can enter through telemetry data supplied by instrumented
  applications, including attributes and log bodies, and through trace context
  and baggage extracted from incoming requests.
- Network input also includes responses from telemetry export destinations and
  requests to the Prometheus exporter's HTTP handler, when enabled.

## Components that matter most / least

- All production packages in this repository are in scope. Prioritize context
  propagation and baggage parsing, SDK processing and resource management, and
  exporters that communicate over the network.
- Tests, test helpers, and examples support exercising these components and are
  lower priority than production code.

## How to exercise it

- Unit tests, integration tests, and executable examples are alongside the code
  in `*_test.go` files. Fuzz tests include those in `trace/trace_fuzz_test.go`.
- Run `make test` to test all modules. See [CONTRIBUTING.md](CONTRIBUTING.md) for
  development and verification instructions.
- See [.oss-scanner/README.md](.oss-scanner/README.md) for the prepared offline
  environment. The image includes the shared security policy at
  `/opt/oss-scanner/SECURITY.md`; this threat model's overrides still apply.

## How you rate severity

Follow the severity guidance and report requirements in the OpenTelemetry
Security Policy, assessing impact on both the SDK and the application using it.

### Availability override

For all artifacts in this repository, the following are considered security
vulnerabilities:

- A denial-of-service attack by a properly authenticated client.
- An availability-related attack on an endpoint by a properly authenticated
  client.

This explicitly overrides the two authenticated-client exclusions in the
[Security Policy's security model](https://github.com/open-telemetry/.github/blob/main/SECURITY.md#security-model).
Authentication does not put these attacks out of scope. The policy's remaining
guidance, including its availability assessment and private reporting process,
continues to apply.
