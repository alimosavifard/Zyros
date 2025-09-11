// controllers/likeController.go

package controllers

import (
    "github.com/alimosavifard/zyros-backend/services"
    "github.com/alimosavifard/zyros-backend/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
)

type LikeController struct {
    service *services.LikeService
}

func NewLikeController(service *services.LikeService) *LikeController {
    return &LikeController{service: service}
}

func (c *LikeController) LikePost(ctx *gin.Context) {
    postID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        utils.SendError(ctx, http.StatusBadRequest, "Invalid post ID", nil)
        return
    }

    userID, exists := ctx.Get("userID")
    if !exists {
        utils.SendError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
        return
    }

    if err := c.service.LikePost(userID.(uint), uint(postID)); err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to like post", err)
        return
    }
    utils.SendSuccess(ctx, "Post liked successfully", nil, nil)
}

func (c *LikeController) UnlikePost(ctx *gin.Context) {
    // Similar to LikePost, but calls UnlikePost service method
}