// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpconf // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/otlpconf"

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/retry"
)

// Default values.
var (
	defaultEndpoint = "localhost:4318"
	defaultPath     = "/v1/logs"
	defaultTimeout  = 10 * time.Second
	defaultRetryCfg = retry.DefaultConfig
)

// Environment variable keys.
var (
	envEndpoint = []string{
		"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	}
	envInsecure = envEndpoint

	// Split because these are parsed differently.
	envPathSignal = []string{"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT"}
	envPathOTLP   = []string{"OTEL_EXPORTER_OTLP_ENDPOINT"}

	envHeaders = []string{
		"OTEL_EXPORTER_OTLP_LOGS_HEADERS",
		"OTEL_EXPORTER_OTLP_HEADERS",
	}

	envCompression = []string{
		"OTEL_EXPORTER_OTLP_LOGS_COMPRESSION",
		"OTEL_EXPORTER_OTLP_COMPRESSION",
	}

	envTimeout = []string{
		"OTEL_EXPORTER_OTLP_LOGS_TIMEOUT",
		"OTEL_EXPORTER_OTLP_TIMEOUT",
	}

	envTLSCert = []string{
		"OTEL_EXPORTER_OTLP_LOGS_CERTIFICATE",
		"OTEL_EXPORTER_OTLP_CERTIFICATE",
	}
	envTLSClient = []struct {
		Certificate string
		Key         string
	}{
		{
			"OTEL_EXPORTER_OTLP_LOGS_CLIENT_CERTIFICATE",
			"OTEL_EXPORTER_OTLP_LOGS_CLIENT_KEY",
		},
		{
			"OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE",
			"OTEL_EXPORTER_OTLP_CLIENT_KEY",
		},
	}
)

// Setting is a configuration setting value.
type Setting[T any] struct {
	Value T
	Set   bool
}

// NewSetting returns a new setting with the value set.
func NewSetting[T any](value T) Setting[T] {
	return Setting[T]{Value: value, Set: true}
}

// Resolver returns an updated setting after applying an resolution operation.
type Resolver[T any] func(Setting[T]) Setting[T]

// Resolve returns a resolved version of s.
//
// It will apply all the passed fn in the order provided, chaining together the
// return setting to the next input. The setting s is used as the initial
// argument to the first fn.
//
// Each fn needs to validate if it should apply given the Set state of the
// setting. This will not perform any checks on the set state when chaining
// function.
func (s Setting[T]) Resolve(fn ...Resolver[T]) Setting[T] {
	for _, f := range fn {
		s = f(s)
	}
	return s
}

// Compression describes the compression used for exported payloads.
type Compression int

const (
	// NoCompression represents that no compression should be used.
	NoCompression Compression = iota
	// GzipCompression represents that gzip compression should be used.
	GzipCompression
)

// Option applies an option to the Exporter.
type Option interface {
	ApplyOption(Config) Config
}

type fnOpt func(Config) Config

func (f fnOpt) ApplyOption(c Config) Config { return f(c) }

type Config struct {
	Endpoint    Setting[string]
	Path        Setting[string]
	Insecure    Setting[bool]
	TLSCfg      Setting[*tls.Config]
	Headers     Setting[map[string]string]
	Compression Setting[Compression]
	Timeout     Setting[time.Duration]
	RetryCfg    Setting[retry.Config]
}

func LoadConfig(c Config) Config {
	c.Endpoint = c.Endpoint.Resolve(
		GetEnv[string](envEndpoint, convEndpoint),
		fallback[string](defaultEndpoint),
	)
	c.Path = c.Path.Resolve(
		GetEnv[string](envPathSignal, convPathExact),
		GetEnv[string](envPathOTLP, convPath),
		fallback[string](defaultPath),
	)
	c.Insecure = c.Insecure.Resolve(
		GetEnv[bool](envInsecure, convInsecure),
	)
	c.TLSCfg = c.TLSCfg.Resolve(
		loadEnvTLS[*tls.Config](),
	)
	c.Headers = c.Headers.Resolve(
		GetEnv[map[string]string](envHeaders, convHeaders),
	)
	c.Compression = c.Compression.Resolve(
		GetEnv[Compression](envCompression, convCompression),
	)
	c.Timeout = c.Timeout.Resolve(
		GetEnv[time.Duration](envTimeout, convDuration),
		fallback[time.Duration](defaultTimeout),
	)
	c.RetryCfg = c.RetryCfg.Resolve(
		fallback[retry.Config](defaultRetryCfg),
	)

	return c
}

// GetEnv returns a Resolver that will apply an environment variable value
// associated with the first set key to a setting value. The conv function is
// used to convert between the environment variable value and the setting type.
//
// If the input setting to the Resolver is set, the environment variable will
// not be applied.
//
// Any error returned from conv is sent to the OTel ErrorHandler and the
// setting will not be updated.
func GetEnv[T any](keys []string, conv func(string) (T, error)) Resolver[T] {
	return func(s Setting[T]) Setting[T] {
		if s.Set {
			// Passed, valid, options have precedence.
			return s
		}

		for _, key := range keys {
			if vStr := os.Getenv(key); vStr != "" {
				v, err := conv(vStr)
				if err == nil {
					s.Value = v
					s.Set = true
					break
				}
				otel.Handle(fmt.Errorf("invalid %s value %s: %w", key, vStr, err))
			}
		}
		return s
	}
}

// convEndpoint converts s from a URL string to an endpoint if s is a valid
// URL. Otherwise, "" and an error are returned.
func convEndpoint(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// convPathExact converts s from a URL string to the exact path if s is a valid
// URL. Otherwise, "" and an error are returned.
//
// If the path contained in s is empty, "/" is returned.
func convPathExact(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if u.Path == "" {
		return "/", nil
	}
	return u.Path, nil
}

// convPath converts s from a URL string to an OTLP endpoint path if s is a
// valid URL. Otherwise, "" and an error are returned.
func convPath(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Path + "/v1/logs", nil
}

// convInsecure parses s as a URL string and returns if the connection should
// use client transport security or not. If s is an invalid URL, false and an
// error are returned.
func convInsecure(s string) (bool, error) {
	u, err := url.Parse(s)
	if err != nil {
		return false, err
	}
	return u.Scheme != "https", nil
}

// convHeaders converts the OTel environment variable header value s into a
// mapping of header key to value. If s is invalid a partial result and error
// are returned.
func convHeaders(s string) (map[string]string, error) {
	out := make(map[string]string)
	var err error
	for _, header := range strings.Split(s, ",") {
		rawKey, rawVal, found := strings.Cut(header, "=")
		if !found {
			err = errors.Join(err, fmt.Errorf("invalid header: %s", header))
			continue
		}

		escKey, e := url.PathUnescape(rawKey)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid header key: %s", rawKey))
			continue
		}
		key := strings.TrimSpace(escKey)

		escVal, e := url.PathUnescape(rawVal)
		if e != nil {
			err = errors.Join(err, fmt.Errorf("invalid header value: %s", rawVal))
			continue
		}
		val := strings.TrimSpace(escVal)

		out[key] = val
	}
	return out, err
}

// convCompression returns the parsed compression encoded in s. NoCompression
// and an errors are returned if s is unknown.
func convCompression(s string) (Compression, error) {
	switch s {
	case "gzip":
		return GzipCompression, nil
	case "none", "":
		return NoCompression, nil
	}
	return NoCompression, fmt.Errorf("unknown compression: %s", s)
}

// convDuration converts s into a duration of milliseconds. If s does not
// contain an integer, 0 and an error are returned.
func convDuration(s string) (time.Duration, error) {
	d, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	// OTel durations are defined in milliseconds.
	return time.Duration(d) * time.Millisecond, nil
}

// fallback returns a resolve that will set a setting value to val if it is not
// already set.
//
// This is usually passed at the end of a resolver chain to ensure a default is
// applied if the setting has not already been set.
func fallback[T any](val T) Resolver[T] {
	return func(s Setting[T]) Setting[T] {
		if !s.Set {
			s.Value = val
			s.Set = true
		}
		return s
	}
}

// loadEnvTLS returns a resolver that loads a *tls.Config from files defeind by
// the OTLP TLS environment variables. This will load both the rootCAs and
// certificates used for mTLS.
//
// If the filepath defined is invalid or does not contain valid TLS files, an
// error is passed to the OTel ErrorHandler and no TLS configuration is
// provided.
func loadEnvTLS[T *tls.Config]() Resolver[T] {
	return func(s Setting[T]) Setting[T] {
		if s.Set {
			// Passed, valid, options have precedence.
			return s
		}

		var rootCAs *x509.CertPool
		var err error
		for _, key := range envTLSCert {
			if v := os.Getenv(key); v != "" {
				rootCAs, err = loadCertPool(v)
				break
			}
		}

		var certs []tls.Certificate
		for _, pair := range envTLSClient {
			cert := os.Getenv(pair.Certificate)
			key := os.Getenv(pair.Key)
			if cert != "" && key != "" {
				var e error
				certs, e = loadCertificates(cert, key)
				err = errors.Join(err, e)
				break
			}
		}

		if err != nil {
			err = fmt.Errorf("failed to load TLS: %w", err)
			otel.Handle(err)
		} else if rootCAs != nil || certs != nil {
			s.Set = true
			s.Value = &tls.Config{RootCAs: rootCAs, Certificates: certs}
		}
		return s
	}
}

// readFile is used for testing.
var readFile = os.ReadFile

// loadCertPool loads and returns the *x509.CertPool found at path if it exists
// and is valid. Otherwise, nil and an error is returned.
func loadCertPool(path string) (*x509.CertPool, error) {
	b, err := readFile(path)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if ok := cp.AppendCertsFromPEM(b); !ok {
		return nil, errors.New("certificate not added")
	}
	return cp, nil
}

// loadCertificates loads and returns the tls.Certificate found at path if it
// exists and is valid. Otherwise, nil and an error is returned.
func loadCertificates(certPath, keyPath string) ([]tls.Certificate, error) {
	cert, err := readFile(certPath)
	if err != nil {
		return nil, err
	}
	key, err := readFile(keyPath)
	if err != nil {
		return nil, err
	}
	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return []tls.Certificate{crt}, nil
}
