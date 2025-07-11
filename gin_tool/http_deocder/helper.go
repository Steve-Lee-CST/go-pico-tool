package http_decoder

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type helper struct{}

var Helper = &helper{}

func (h *helper) SetHttpRequest(c *gin.Context, config Config, req *HttpRequest) {
	c.Set(config.HttpRequestKey, req)
}

func (h *helper) GetHttpRequest(c *gin.Context, config Config) *HttpRequest {
	if req, exists := c.Get(config.HttpRequestKey); exists {
		if httpReq, ok := req.(*HttpRequest); ok {
			return httpReq
		}
	}
	return nil
}

func (h *helper) SetHttpResponse(c *gin.Context, config Config, resp *HttpResponse) {
	c.Set(config.HttpResponseKey, resp)
}

func (h *helper) GetHttpResponse(c *gin.Context, config Config) *HttpResponse {
	if resp, exists := c.Get(config.HttpResponseKey); exists {
		if httpResp, ok := resp.(*HttpResponse); ok {
			return httpResp
		}
	}
	return nil
}

func (h *helper) SetResponseWriter(c *gin.Context, config Config, writer IWrappedResponseWriter) {
	c.Set(config.HttpResponseWriterKey, writer)
}

func (h *helper) GetResponseWriter(c *gin.Context, config Config) IWrappedResponseWriter {
	if writer, exists := c.Get(config.HttpResponseWriterKey); exists {
		if wrappedWriter, ok := writer.(IWrappedResponseWriter); ok {
			return wrappedWriter
		}
	}
	return nil
}

func (h *helper) DecodeRequest(c *gin.Context, config Config) *HttpRequest {
	req := HttpRequest{
		Method:   c.Request.Method,
		Protocol: c.Request.Proto,
		Host:     c.Request.Host,
		URL:      c.Request.URL,
		FullPath: c.FullPath(),
		IP:       c.ClientIP(),

		Header:    c.Request.Header.Clone(),
		Cookies:   c.Request.Cookies(),
		UserAgent: c.Request.UserAgent(),

		PathParams:  c.Params,
		QueryParams: c.Request.URL.Query(),
		ContentType: c.ContentType(),

		ReceivedTime: time.Now(),
	}

	// read raw body and write back
	body, _ := c.GetRawData()
	if len(body) > 0 {
		req.RawBody = body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// parse json
	if strings.Contains(req.ContentType, "application/json") && len(body) > 0 {
		_ = json.Unmarshal(body, &req.JsonBody)
	}

	// parse form data
	if strings.Contains(req.ContentType, "multipart/form-data") ||
		strings.Contains(req.ContentType, "application/x-www-form-urlencoded") {
		if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
			req.FormValues = c.Request.PostForm
			if c.Request.MultipartForm != nil {
				for _, files := range c.Request.MultipartForm.File {
					for _, file := range files {
						req.FormFiles = append(req.FormFiles, file.Filename)
					}
				}
			}
		}
	}

	return &req
}

func (h *helper) DecodeResponse(
	c *gin.Context, config Config,
) *HttpResponse {
	writer := h.GetResponseWriter(c, config)
	if writer == nil {
		return nil
	}
	resp := HttpResponse{
		Status:      c.Writer.Status(),
		StatusText:  http.StatusText(c.Writer.Status()),
		Size:        c.Writer.Size(),
		ContentType: c.Writer.Header().Get("Content-Type"),

		ResponseTime: time.Now(),
	}

	resp.Header = c.Writer.Header()
	resp.Cookies = c.Request.Cookies()
	resp.Body = writer.GetBodyBytes()

	if strings.Contains(resp.ContentType, "application/json") {
		json.Unmarshal(writer.GetBodyBytes(), &resp.JsonBody)
	}

	return &resp
}
