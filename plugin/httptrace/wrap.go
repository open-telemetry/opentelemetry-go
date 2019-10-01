package httptrace

import (
	"io"
	"net/http"
)

// wrapBody returns the wrapped version of the original io.ReadCloser.
// If the original io.ReadCloser was also an io.Writer the returned
// io.ReadCloser is also an io.Writer
func wrapBody(wrapper io.ReadCloser, original io.ReadCloser) io.ReadCloser {
	w, isWriter := original.(io.Writer)

	if !isWriter {
		return wrapper
	}

	return struct {
		io.ReadCloser
		io.Writer
	}{wrapper, w}
}

var _ io.ReadCloser = &bodyWrapper{}

// bodyWrapper wraps a http.Request.Body (an io.ReadCloser) to track the number
// of bytes read and the last error
type bodyWrapper struct {
	rc io.ReadCloser

	read int64
	err  error
}

func (w *bodyWrapper) Read(b []byte) (int, error) {
	n, err := w.rc.Read(b)
	w.read += int64(n)
	w.err = err
	return n, err
}

func (w *bodyWrapper) Close() error {
	return w.rc.Close()
}

var _ http.ResponseWriter = &respWriterWrapper{}

// respWriterWrapper wraps a http.ResponseWriter in order to track the number of
// bytes written, the last error, and to catch the returned statusCode
// TODO: The wrapped http.ResponseWriter doesn't implement any of the optional
// types (http.Hijacker, http.Pusher, http.CloseNotifier, http.Flusher, etc)
// that may be useful when using it in real life situations.
type respWriterWrapper struct {
	w http.ResponseWriter

	written     int64
	statusCode  int
	err         error
	wroteHeader bool
}

func (w *respWriterWrapper) Header() http.Header {
	return w.w.Header()
}

func (w *respWriterWrapper) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
		w.wroteHeader = true
	}
	n, err := w.w.Write(p)
	w.written += int64(n)
	w.err = err
	return n, err
}

func (w *respWriterWrapper) WriteHeader(statusCode int) {
	w.wroteHeader = true
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}
