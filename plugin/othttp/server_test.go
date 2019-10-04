package othttp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/propagation"

	mocktrace "go.opentelemetry.io/internal/trace"
)

func TestBasics(t *testing.T) {
	rr := httptest.NewRecorder()

	var id uint64
	tracer := mocktrace.MockTracer{StartSpanId: &id}

	h := NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := io.WriteString(w, "hello world"); err != nil {
				t.Fatal(err)
			}
		}), "test_handler",
		WithTracer(&tracer))

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rr, r)
	if got, expected := rr.Result().StatusCode, http.StatusOK; got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	if got := rr.Header().Get(propagation.TraceparentHeader); got == "" {
		t.Fatal("expected non empty trace header")
	}
	if got, expected := id, uint64(1); got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	d, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := string(d), "hello world"; got != expected {
		t.Fatalf("got %q, expected %q", got, expected)
	}
}
