package request_id

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type helper struct{}

var Helper = helper{}

func (h helper) GetRequestIDFromRequest(c *gin.Context, config Config) (string, bool) {
	requestIDs, exists := c.Request.Header[http.CanonicalHeaderKey(
		config.HeaderKey,
	)]
	if !exists || len(requestIDs) == 0 {
		return "", false
	}
	// return first not empty requestID
	for _, requestID := range requestIDs {
		if requestID != "" {
			return requestID, true
		}
	}
	return "", false
}

func (h helper) GetRequestIDFromResponse(c *gin.Context, config Config) (string, bool) {
	requestIDs, exists := c.Writer.Header()[http.CanonicalHeaderKey(config.HeaderKey)]
	if !exists || len(requestIDs) == 0 {
		return "", false
	}
	// return first not empty requestID
	for _, requestID := range requestIDs {
		if requestID != "" {
			return requestID, true
		}
	}
	return "", false
}

func (h helper) SetRequestIDToRequest(c *gin.Context, config Config, requestID string) {
	if requestID == "" {
		return
	}
	// set requestID to request header
	c.Request.Header.Set(config.HeaderKey, requestID)
}

func (h helper) SetRequestIDToResponse(c *gin.Context, config Config, requestID string) {
	if requestID == "" {
		return
	}
	// set requestID to response header
	c.Writer.Header().Set(config.HeaderKey, requestID)
}
