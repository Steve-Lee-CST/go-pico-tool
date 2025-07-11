package http_decoder

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// HttpRequest some most common fields of HTTP request
type HttpRequest struct {
	Method   string   `json:"method"`
	Protocol string   `json:"protocol"`
	Host     string   `json:"host"`
	URL      *url.URL `json:"url"`
	Path     string   `json:"path"`
	FullPath string   `json:"full_path"`
	IP       string   `json:"ip"`

	Header    http.Header    `json:"headers"`
	Cookies   []*http.Cookie `json:"cookies"`
	UserAgent string         `json:"user_agent"`

	// Path params
	PathParams gin.Params `json:"path_params,omitempty"`
	// GET query params
	QueryParams url.Values `json:"query_params,omitempty"`

	// POST content-type
	ContentType string `json:"content_type"`
	// POST body
	RawBody  []byte         `json:"raw_body,omitempty"`
	JsonBody map[string]any `json:"json_body,omitempty"`
	// POST form data
	FormValues url.Values `json:"form_values,omitempty"`
	// POST form files
	FormFiles []string `json:"form_files,omitempty"`

	// time
	ReceivedTime time.Time `json:"received_time"`
}

// HttpResponse some most common fields of HTTP response
type HttpResponse struct {
	// response info
	Status      int    `json:"status"`
	StatusText  string `json:"status_text"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`

	// Header And Cookies
	Header  http.Header    `json:"headers,omitempty"`
	Cookies []*http.Cookie `json:"cookies,omitempty"`

	// response body
	Body     []byte         `json:"body,omitempty"`
	JsonBody map[string]any `json:"json_body,omitempty"`

	// error info
	Error      string `json:"error,omitempty"`
	StackTrace string `json:"stack_trace,omitempty"`

	ResponseTime time.Time `json:"response_time"`
}

// Check if WrappedResponseWriter implements IWrappedResponseWriter and gin.ResponseWriter
var (
	_ IWrappedResponseWriter = &WrappedResponseWriter{}
	_ gin.ResponseWriter     = &WrappedResponseWriter{}
)

/*
IWrappedResponseWriter And WrappedResponseWriter:
    1. IWrappedResponseWriter is a interface to wrap gin.ResponseWriter,
        which can be used to get response body.
    2. WrappedResponseWriter is a default implementation of IWrappedResponseWriter,
        it will copy response body to a bytes.Buffer.
*/
// IWrappedResponseWriter is an interface that extends gin.ResponseWriter
type IWrappedResponseWriter interface {
	gin.ResponseWriter
	GetBodyBytes() []byte
}

type WrappedResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewWrappedResponseWriter(w gin.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r *WrappedResponseWriter) GetBodyBytes() []byte {
	return r.body.Bytes()
}

func (r *WrappedResponseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)                  // capture response body
	return r.ResponseWriter.Write(b) // write to original response
}
