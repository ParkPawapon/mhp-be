package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ParkPawapon/mhp-be/internal/constants"
)

type Meta struct {
	RequestID string `json:"request_id"`
	Page      int    `json:"page,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

type SuccessResponse struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type ErrorResponse struct {
	Error ErrorObject `json:"error"`
	Meta  Meta        `json:"meta"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Data: data,
		Meta: Meta{RequestID: requestID(c)},
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Data: data,
		Meta: Meta{RequestID: requestID(c)},
	})
}

func Respond(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Data: data,
		Meta: Meta{RequestID: requestID(c)},
	})
}

func Fail(c *gin.Context, err error) {
	httpErr := MapError(err)
	c.JSON(httpErr.Status, ErrorResponse{
		Error: ErrorObject{
			Code:    httpErr.Code,
			Message: httpErr.Message,
			Details: httpErr.Details,
		},
		Meta: Meta{RequestID: requestID(c)},
	})
}

func requestID(c *gin.Context) string {
	if v, ok := c.Get(constants.RequestIDKey); ok {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
