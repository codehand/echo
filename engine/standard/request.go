package standard

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/webx-top/echo/engine"
)

var defaultMaxRequestBodySize int64 = 32 << 20 // 32 MB

type (
	Request struct {
		config  *engine.Config
		request *http.Request
		url     engine.URL
		header  engine.Header
		value   *Value
	}
)

func NewRequest(r *http.Request) *Request {
	req := &Request{
		request: r,
		url:     &URL{url: r.URL},
		header:  &Header{r.Header},
	}
	req.value = NewValue(req)
	return req
}

func (r *Request) Host() string {
	return r.request.Host
}

func (r *Request) URL() engine.URL {
	return r.url
}

func (r *Request) Header() engine.Header {
	return r.header
}

func (r *Request) Proto() string {
	return r.request.Proto
}

//
// func ProtoMajor() int {
// 	return r.request.ProtoMajor()
// }
//
// func ProtoMinor() int {
// 	return r.request.ProtoMinor()
// }

func (r *Request) RemoteAddress() string {
	return r.request.RemoteAddr
}

func (r *Request) Method() string {
	return r.request.Method
}

func (r *Request) SetMethod(method string) {
	r.request.Method = method
}

func (r *Request) URI() string {
	return r.request.RequestURI
}

func (r *Request) Body() io.ReadCloser {
	return r.request.Body
}

// SetBody implements `engine.Request#SetBody` function.
func (r *Request) SetBody(reader io.Reader) {
	r.request.Body = ioutil.NopCloser(reader)
}

func (r *Request) FormValue(name string) string {
	return r.request.FormValue(name)
}

func (r *Request) Form() engine.URLValuer {
	return r.value
}

func (r *Request) PostForm() engine.URLValuer {
	return r.value.postArgs
}

func (r *Request) MultipartForm() *multipart.Form {
	if r.request.MultipartForm == nil {
		maxMemory := defaultMaxRequestBodySize
		if r.config.MaxRequestBodySize != 0 {
			maxMemory = int64(r.config.MaxRequestBodySize)
		}
		r.request.ParseMultipartForm(maxMemory)
	}
	return r.request.MultipartForm
}

func (r *Request) IsTLS() bool {
	return r.request.TLS != nil
}

func (r *Request) Cookie(key string) string {
	if cookie, err := r.request.Cookie(key); err == nil {
		return cookie.Value
	}
	return ``
}

func (r *Request) Referer() string {
	return r.request.Referer()
}

func (r *Request) UserAgent() string {
	return r.request.UserAgent()
}

func (r *Request) Object() interface{} {
	return r.request
}

func (r *Request) reset(req *http.Request, h engine.Header, u engine.URL) {
	r.request = req
	r.header = h
	r.url = u
	r.value = NewValue(r)
}

func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	file, fileHeader, err := r.request.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	return file, fileHeader, err
}

// Size implements `engine.Request#ContentLength` function.
func (r *Request) Size() int64 {
	return r.request.ContentLength
}

func (r *Request) Scheme() string {
	if r.IsTLS() {
		return `https`
	}
	return `http`
}

func (r *Request) BasicAuth() (username, password string, ok bool) {
	return r.request.BasicAuth()
}
