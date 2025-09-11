package controllers

import (
    "github.com/alimosavifard/zyros-backend/config"
    "github.com/alimosavifard/zyros-backend/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

func HealthCheck(ctx *gin.Context) {
    dbStatus := "up"
    if db, err := config.ConnectDB().DB(); err != nil {
        dbStatus = "down"
    } else {
        if err := db.Ping(); err != nil {
            dbStatus = "down"
        }
    }

    redisStatus := "up"
    if err := config.ConnectRedis().Ping(ctx).Err(); err != nil {
        redisStatus = "down"
    }

    data := gin.H{
        "database": dbStatus,
        "redis":    redisStatus,
    }

    utils.SendSuccess(ctx, "Service is healthy", data, nil)
}