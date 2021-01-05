// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlphttp_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	mathrand "math/rand"
	"net"
	"time"
)

type mathRandReader struct{}

func (mathRandReader) Read(p []byte) (n int, err error) {
	return mathrand.Read(p)
}

var randReader mathRandReader

type pemCertificate struct {
	Certificate []byte
	PrivateKey  []byte
}

// Based on https://golang.org/src/crypto/tls/generate_cert.go,
// simplified and weakened.
func generateWeakCertificate() (*pemCertificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), randReader)
	if err != nil {
		return nil, err
	}
	keyUsage := x509.KeyUsageDigitalSignature
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := cryptorand.Int(randReader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"otel-go"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv6loopback, net.IPv4(127, 0, 0, 1)},
	}
	derBytes, err := x509.CreateCertificate(randReader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}
	certificateBuffer := new(bytes.Buffer)
	if err := pem.Encode(certificateBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, err
	}
	privDERBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	privBuffer := new(bytes.Buffer)
	if err := pem.Encode(privBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privDERBytes}); err != nil {
		return nil, err
	}
	return &pemCertificate{
		Certificate: certificateBuffer.Bytes(),
		PrivateKey:  privBuffer.Bytes(),
	}, nil
}
