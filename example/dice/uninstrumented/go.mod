module go.opentelemetry.io/otel/example/dice/uninstrumented

go 1.22

replace go.opentelemetry.io/otel/exporters/stdout/stdouttrace => ./../../../exporters/stdout/stdouttrace

replace go.opentelemetry.io/otel/exporters/stdout/stdoutmetric => ./../../../exporters/stdout/stdoutmetric

replace go.opentelemetry.io/otel => ./../../..

replace go.opentelemetry.io/otel/trace => ./../../../trace

replace go.opentelemetry.io/otel/metric => ./../../../metric

replace go.opentelemetry.io/otel/sdk/metric => ./../../../sdk/metric

replace go.opentelemetry.io/otel/sdk => ./../../../sdk

replace go.opentelemetry.io/otel/exporters/stdout/stdoutlog => ./../../../exporters/stdout/stdoutlog

replace go.opentelemetry.io/otel/log => ./../../../log

replace go.opentelemetry.io/otel/sdk/log => ./../../../sdk/log
