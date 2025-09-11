package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type StandardResponse struct {
    Data    interface{} `json:"data,omitempty"`
    Message string      `json:"message,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SendSuccess(ctx *gin.Context, message string, data interface{}, meta interface{}) {
    response := StandardResponse{
        Data:    data,
        Message: message,
        Meta:    meta,
    }
    ctx.JSON(http.StatusOK, response)
}

func SendError(ctx *gin.Context, statusCode int, message string, err error) {
    if err != nil {
        InitLogger().Error().Err(err).Msg(message)
    }

    response := StandardResponse{
        Error: message,
    }
    ctx.JSON(statusCode, response)
}