# Baggage

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/baggage)](https://pkg.go.dev/go.opentelemetry.io/otel/baggage)

Use [`ContextWithBaggage`](https://pkg.go.dev/go.opentelemetry.io/otel/baggage#ContextWithBaggage) to store baggage in a local context, and [`NewMemberRaw`](https://pkg.go.dev/go.opentelemetry.io/otel/baggage#NewMemberRaw) to create a member from raw key and value strings. To inject or extract baggage across process boundaries, see [`propagation.Baggage`](https://pkg.go.dev/go.opentelemetry.io/otel/propagation#Baggage).
