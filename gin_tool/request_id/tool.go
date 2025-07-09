package request_id

import (
	"net/http"

	"github.com/Steve-Lee-CST/go-pico-tool/gin_tool/common"
	"github.com/Steve-Lee-CST/go-pico-tool/pico_tool/id_generator"
	"github.com/gin-gonic/gin"
)

type RequestIDTool struct {
	config      Config
	idGenerator *id_generator.IDGenerator
}

func NewRequestIDTool(config Config) *RequestIDTool {
	idGenerator := id_generator.NewIDGenerator(config.IDGeneratorConfig)
	return &RequestIDTool{
		config:      config,
		idGenerator: idGenerator,
	}
}

func (tool *RequestIDTool) GenerateRequestID() string {
	return tool.idGenerator.Generate()
}

func (tool *RequestIDTool) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request id from request-header
		requestID, exist := Helper.GetRequestIDFromRequest(c, tool.config)
		if !exist {
			requestID = tool.GenerateRequestID()
			Helper.SetRequestIDToRequest(c, tool.config, requestID)
		}
		// Set the request ID in the response-header
		Helper.SetRequestIDToResponse(c, tool.config, requestID)
	}
}

func (tool *RequestIDTool) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := tool.GenerateRequestID()
		// return requestID from 2 ways: response & response-header
		Helper.SetRequestIDToResponse(c, tool.config, requestID)
		c.JSON(http.StatusOK, common.SuccessResponse(&requestID))
	}
}
