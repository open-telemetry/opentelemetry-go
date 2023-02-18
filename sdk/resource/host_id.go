package resource

import (
	"context"
	"os"
	"os/exec"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type hostIDProvider func() (string, error)

var defaultHostIDProvider hostIDProvider = platformHostIDReader.read

var hostID = defaultHostIDProvider

func setDefaultHostIDProvider() {
	setHostIDProvider(defaultHostIDProvider)
}

func setHostIDProvider(hostIDProvider hostIDProvider) {
	hostID = hostIDProvider
}

type hostIDReader interface {
	read() (string, error)
}

type fileReader func(string) (string, error)

type commandExecutor func(string, ...string) (string, error)

func readFile(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

func execCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type hostIDDetector struct{}

// Detect returns a *Resource containing the platform specific host id
func (hostIDDetector) Detect(ctx context.Context) (*Resource, error) {
	hostID, err := hostID()
	if err != nil {
		return nil, err
	}

	return NewWithAttributes(
		semconv.SchemaURL,
		semconv.HostID(hostID),
	), nil
}
