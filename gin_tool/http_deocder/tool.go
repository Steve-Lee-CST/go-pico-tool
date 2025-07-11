package http_decoder

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpDecoder struct {
	config Config
}

func NewHttpDecoder(config Config) *HttpDecoder {
	return &HttpDecoder{
		config: config,
	}
}

func (hd *HttpDecoder) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decode request
		req := Helper.DecodeRequest(c, hd.config)
		if req != nil {
			Helper.SetHttpRequest(c, hd.config, req)
		}

		// Create a wrapped response writer
		writer := NewWrappedResponseWriter(c.Writer)
		Helper.SetResponseWriter(c, hd.config, writer)

		c.Next()

		// Decode response
		resp := Helper.DecodeResponse(c, hd.config)
		if resp != nil {
			Helper.SetHttpResponse(c, hd.config, resp)
		}
	}
}

func (hd *HttpDecoder) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpRequest := Helper.GetHttpRequest(c, hd.config)
		if httpRequest == nil {
			c.JSON(400, gin.H{"error": "No HTTP request found"})
			return
		}
		writer := NewWrappedResponseWriter(c.Writer)
		Helper.SetResponseWriter(c, hd.config, writer)
		c.JSON(http.StatusOK, httpRequest)
	}
}
