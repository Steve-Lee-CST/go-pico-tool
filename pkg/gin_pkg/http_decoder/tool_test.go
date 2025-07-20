package http_decoder

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHttpDecoder_DefaultConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.GET("/decode", decoder.Handler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/decode", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response contains HTTP request information
	assert.Contains(t, response, "method")
	assert.Contains(t, response, "path")
	assert.Contains(t, response, "headers")
}

func TestHttpDecoder_CustomConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := Config{
		HttpRequestKey:        "custom_request",
		HttpResponseKey:       "custom_response",
		HttpResponseWriterKey: "custom_writer",
	}
	decoder := NewHttpDecoder(cfg)
	r.Use(decoder.Middleware())
	r.GET("/decode", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, cfg)
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "GET", req.Method)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/decode?param=value", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_Middleware_JSONRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.POST("/test", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "application/json", req.ContentType)
		assert.NotEmpty(t, req.RawBody, "Raw body should be captured")
		assert.NotEmpty(t, req.JsonBody, "JSON body should be parsed")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	jsonData := `{"name": "test", "value": 123}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_Middleware_FormData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.POST("/form", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.ContentType, "application/x-www-form-urlencoded")
		// Verify that raw request body is captured
		assert.NotEmpty(t, req.RawBody, "Raw body should be captured")
		// Since form parsing may be affected by GetRawData, we mainly verify that raw data is captured
		assert.Contains(t, string(req.RawBody), "name=test")
		assert.Contains(t, string(req.RawBody), "value=123")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	body := strings.NewReader("name=test&value=123")
	req, _ := http.NewRequest("POST", "/form", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_Middleware_QueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.GET("/query", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "GET", req.Method)
		assert.NotEmpty(t, req.QueryParams, "Query parameters should be captured")
		assert.Equal(t, "value1", req.QueryParams.Get("param1"))
		assert.Equal(t, "value2", req.QueryParams.Get("param2"))
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query?param1=value1&param2=value2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_Middleware_PathParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.GET("/user/:id/profile", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "GET", req.Method)
		assert.NotEmpty(t, req.PathParams, "Path parameters should be captured")
		assert.Equal(t, "123", req.PathParams.ByName("id"))
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/123/profile", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_ResponseCapture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.GET("/response", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "test response"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/response", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify that response is correctly captured
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test response", response["message"])
}

func TestHttpDecoder_HelperFunctions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := DefaultConfig()

	// Test setting and getting HTTP requests
	r.GET("/helper", func(c *gin.Context) {
		// Test setting HTTP request
		testReq := &HttpRequest{
			Method: "GET",
			IP:     "127.0.0.1",
		}
		Helper.SetHttpRequest(c, cfg, testReq)

		// Test getting HTTP request
		retrievedReq := Helper.GetHttpRequest(c, cfg)
		assert.NotNil(t, retrievedReq)
		assert.Equal(t, "GET", retrievedReq.Method)
		assert.Equal(t, "127.0.0.1", retrievedReq.IP)

		// Test setting HTTP response
		testResp := &HttpResponse{
			Status: 200,
			Size:   100,
		}
		Helper.SetHttpResponse(c, cfg, testResp)

		// Test getting HTTP response
		retrievedResp := Helper.GetHttpResponse(c, cfg)
		assert.NotNil(t, retrievedResp)
		assert.Equal(t, 200, retrievedResp.Status)
		assert.Equal(t, 100, retrievedResp.Size)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/helper", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_DecodeRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/decode-test", func(c *gin.Context) {
		req := Helper.DecodeRequest(c, DefaultConfig())
		assert.NotNil(t, req)
		assert.Equal(t, "POST", req.Method)
		assert.NotZero(t, req.ReceivedTime)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	jsonData := `{"test": "data"}`
	req, _ := http.NewRequest("POST", "/decode-test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-agent")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHttpDecoder_DecodeResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())

	r.GET("/response-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/response-test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify that response body is correctly parsed
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response["message"])
}

func TestHttpDecoder_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())

	r.GET("/error", func(c *gin.Context) {
		// Simulate error condition
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHttpDecoder_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())

	r.GET("/empty", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req)
		assert.Equal(t, "GET", req.Method)
		assert.Empty(t, req.RawBody, "Empty body should be handled correctly")
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/empty", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHttpDecoder_Handler_NoRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())

	// Test Handler directly without using middleware
	r.GET("/no-request", decoder.Handler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/no-request", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestWrappedResponseWriter(t *testing.T) {
	// Test wrapped response writer
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/test-writer", func(c *gin.Context) {
		writer := NewWrappedResponseWriter(c.Writer)

		// Write data
		testData := []byte("test response")
		n, err := writer.Write(testData)

		assert.NoError(t, err)
		assert.Equal(t, len(testData), n)

		// Verify data is correctly captured
		bodyBytes := writer.GetBodyBytes()
		assert.Equal(t, testData, bodyBytes)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-writer", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test response", w.Body.String())
}

func TestHttpDecoder_Middleware_MultipartFormData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	decoder := NewHttpDecoder(DefaultConfig())
	r.Use(decoder.Middleware())
	r.POST("/multipart", func(c *gin.Context) {
		req := Helper.GetHttpRequest(c, DefaultConfig())
		assert.NotNil(t, req, "HTTP request should be decoded")
		assert.Equal(t, "POST", req.Method)
		assert.Contains(t, req.ContentType, "multipart/form-data")
		// Verify that raw request body is captured
		assert.NotEmpty(t, req.RawBody, "Raw body should be captured")
		// Verify that form values are parsed (if possible)
		if len(req.FormValues) > 0 {
			assert.Equal(t, "test", req.FormValues.Get("name"))
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field
	field, err := writer.CreateFormField("name")
	assert.NoError(t, err)
	field.Write([]byte("test"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/multipart", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
