package request_id

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Steve-Lee-CST/go-pico-tool/gin_tool/common"
	"github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDTool_DefaultConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tool := NewRequestIDTool(GetDefaultConfig())
	r.Use(tool.Middleware())
	r.GET("/request_id", tool.Handler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/request_id", nil)
	r.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "Request ID should not be empty")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), requestID)
}

func TestRequestIDTool_CustomConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	sep := "_"
	customModifier := func(timestamp int64, microSecond int64, randSegment string) []string {
		return []string{
			"CUSTOM",
			fmt.Sprintf("%d", timestamp),
			randSegment,
		}
	}
	cfg := Config{
		HeaderKey: "X-Custom-Request-ID",
		IDGeneratorConfig: id_generator.Config{
			Separator: &sep,
			Modifier:  customModifier,
		},
	}
	tool := NewRequestIDTool(cfg)
	r.Use(tool.Middleware())
	r.GET("/request_id", tool.Handler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/request_id", nil)
	r.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Custom-Request-ID")
	assert.NotEmpty(t, requestID, "Custom Request ID should not be empty")
	assert.True(t, strings.HasPrefix(requestID, "CUSTOM"), "Custom Request ID should have prefix 'CUSTOM'")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), requestID)
}

func TestRequestIDTool_ReuseRequestIDFromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tool := NewRequestIDTool(GetDefaultConfig())
	r.Use(tool.Middleware())
	r.GET("/request_id", func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = tool.GenerateRequestID()
		}
		Helper.SetRequestIDToResponse(c, tool.config, requestID)
		c.JSON(http.StatusOK, common.SuccessResponse(&requestID))
	})

	w := httptest.NewRecorder()
	customID := "test-reuse-id"
	req, _ := http.NewRequest("GET", "/request_id", nil)
	req.Header.Set("X-Request-ID", customID)
	r.ServeHTTP(w, req)

	requestID := w.Header().Get("X-Request-ID")
	assert.Equal(t, customID, requestID, "Should reuse the request ID from header")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), requestID)
}
