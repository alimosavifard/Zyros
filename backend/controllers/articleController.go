package controllers

import (
	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/requests"
	"github.com/alimosavifard/zyros-backend/services"
	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleController struct {
    service *services.PostService
}

func NewArticleController(service *services.PostService) *ArticleController {
    return &ArticleController{service: service}
}

func (ctrl *ArticleController) CreateArticle(c *gin.Context) {
    var req requests.ArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.SendError(c, http.StatusBadRequest, "Invalid input", err)
        return
    }
    
    if err := requests.ValidateStruct(&req); err != nil {
        utils.SendError(c, http.StatusBadRequest, "Validation failed", err)
        return
    }

    userID, exists := c.Get("userID")
    if !exists {
        utils.SendError(c, http.StatusUnauthorized, "Unauthorized", nil)
        return
    }

    article := &models.Post{
        Title:    req.Title,
        Content:  req.Content,
        Type:     "article",
        Lang:     req.Lang,
        UserID:   userID.(uint),
        ImageUrl: req.ImageUrl,
    }

    if err := ctrl.service.CreatePost(c, article); err != nil {
        utils.SendError(c, http.StatusInternalServerError, "Failed to create article", err)
        return
    }

    utils.SendSuccess(c, "Article created successfully", article, nil)
}