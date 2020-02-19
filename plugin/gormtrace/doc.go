//Package gormtrace allows for the wrapping of GORM calls to
//databases with OpenTelemetry tracing spans. You only need to
//create your GORM db client and pass that into otgorm.WithContext
//along with context.Context(). If there is a parent span referenced
//within the context the GORM call will be a child span.
package gormtrace
