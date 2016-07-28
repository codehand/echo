package engine

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/webx-top/echo/logger"
)

type (
	// Engine defines an interface for HTTP server.
	Engine interface {
		SetHandler(Handler)
		SetLogger(logger.Logger)
		Start() error
	}

	// Request defines an interface for HTTP request.
	Request interface {
		// Scheme returns the HTTP protocol scheme, `http` or `https`.
		Scheme() string

		// Host returns HTTP request host. Per RFC 2616, this is either the value of
		// the `Host` header or the host name given in the URL itself.
		Host() string

		// URI returns the unmodified `Request-URI` sent by the client.
		URI() string

		// URL returns `engine.URL`.
		URL() URL

		// Header returns `engine.Header`.
		Header() Header

		// Proto returns the HTTP proto. (HTTP/1.1 etc.)
		Proto() string
		// ProtoMajor() int
		// ProtoMinor() int

		// RemoteAddress returns the client's network address.
		RemoteAddress() string

		// Method returns the request's HTTP function.
		Method() string

		// SetMethod sets the HTTP method of the request.
		SetMethod(string)

		// Body returns request's body.
		Body() io.ReadCloser

		SetBody(io.Reader)

		// FormValue returns the form field value for the provided name.
		FormValue(string) string
		Object() interface{}

		Form() URLValuer
		PostForm() URLValuer

		// MultipartForm returns the multipart form.
		MultipartForm() *multipart.Form

		// IsTLS returns true if HTTP connection is TLS otherwise false.
		IsTLS() bool
		Cookie(string) string
		Referer() string

		// UserAgent returns the client's `User-Agent`.
		UserAgent() string

		// FormFile returns the multipart form file for the provided name.
		FormFile(string) (multipart.File, *multipart.FileHeader, error)

		// ContentLength returns the size of request's body.
		Size() int64

		BasicAuth() (string, string, bool)
	}

	// Response defines an interface for HTTP response.
	Response interface {
		// Header returns `engine.Header`
		Header() Header

		// WriteHeader sends an HTTP response header with status code.
		WriteHeader(int)

		// Write writes the data to the connection as part of an HTTP reply.
		Write(b []byte) (int, error)

		// Status returns the HTTP response status.
		Status() int

		// Size returns the number of bytes written to HTTP response.
		Size() int64

		// Committed returns true if HTTP response header is written, otherwise false.
		Committed() bool

		// SetWriter sets the HTTP response writer.
		SetWriter(io.Writer)

		// Write returns the HTTP response writer.
		Writer() io.Writer
		Object() interface{}

		Hijack(func(net.Conn))
		Body() []byte
		Redirect(string, int)
		NotFound()
		SetCookie(*http.Cookie)
		ServeFile(string)
	}

	// Header defines an interface for HTTP header.
	Header interface {
		// Add adds the key, value pair to the header. It appends to any existing values
		// associated with key.
		Add(string, string)

		// Del deletes the values associated with key.
		Del(string)

		// Get gets the first value associated with the given key. If there are
		// no values associated with the key, Get returns "".
		Get(string) string

		// Set sets the header entries associated with key to the single element value.
		// It replaces any existing values associated with key.
		Set(string, string)

		Object() interface{}
	}

	// URLValuer Wrap url.Values
	URLValuer interface {
		Add(string, string)
		Del(string)
		Get(string) string
		Gets(string) []string
		Set(string, string)
		Encode() string
		All() map[string][]string
		Reset(url.Values)
	}

	// URL defines an interface for HTTP request url.
	URL interface {
		SetPath(string)
		Path() string
		QueryValue(string) string
		QueryValues(string) []string
		Query() url.Values
		RawQuery() string
		Object() interface{}
	}

	// Config defines engine configuration.
	Config struct {
		Address            string        // TCP address to listen on.
		Listener           net.Listener  // Custom `net.Listener`. If set, server accepts connections on it.
		TLSCertfile        string        // TLS certificate file path.
		TLSKeyfile         string        // TLS key file path.
		ReadTimeout        time.Duration // Maximum duration before timing out read of the request.
		WriteTimeout       time.Duration // Maximum duration before timing out write of the response.
		MaxConnsPerIP      int
		MaxRequestsPerConn int
		MaxRequestBodySize int
	}

	// Handler defines an interface to server HTTP requests via `ServeHTTP(Request, Response)`
	// function.
	Handler interface {
		ServeHTTP(Request, Response)
	}

	// HandlerFunc is an adapter to allow the use of `func(Request, Response)` as HTTP handlers.
	HandlerFunc func(Request, Response)
)

// ServeHTTP serves HTTP request.
func (h HandlerFunc) ServeHTTP(req Request, res Response) {
	h(req, res)
}
