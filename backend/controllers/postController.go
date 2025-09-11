package controllers

import (
	"context"
    "fmt"
    "github.com/alimosavifard/zyros-backend/models"
    "github.com/alimosavifard/zyros-backend/requests"
    "github.com/alimosavifard/zyros-backend/services"
    "github.com/alimosavifard/zyros-backend/utils"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

type PostController struct {
    postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
    return &PostController{postService: postService}
}

func (c *PostController) CreatePost(ctx *gin.Context) {
    var req requests.PostRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.SendError(ctx, http.StatusBadRequest, "Invalid input", err)
        return
    }
    
    if err := requests.ValidateStruct(&req); err != nil {
        utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err)
        return
    }

    userID, exists := ctx.Get("userID")
    if !exists {
        utils.SendError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
        return
    }

    post := &models.Post{
        Title:    req.Title,
        Content:  req.Content,
        Type:     req.Type,
        Lang:     req.Lang,
        UserID:   userID.(uint),
        ImageUrl: req.ImageUrl,
    }

    if err := c.postService.CreatePost(ctx, post); err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to create post", err)
        return
    }

    utils.SendSuccess(ctx, "Post created successfully", post, nil)
}


func (c *PostController) GetPosts(ctx *gin.Context) {
    lang := ctx.DefaultQuery("lang", "en")
    postType := ctx.DefaultQuery("type", "post")
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }

    posts, err := c.postService.GetPosts(lang, postType, page, limit)
    if err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve posts", err)
        return
    }

    // This is the important part, it sends a response with the "posts" key inside "data"
    utils.SendSuccess(ctx, "Posts retrieved successfully", gin.H{"posts": posts}, nil)
}


// متد جدید برای دریافت یک پست با شناسه
func (c *PostController) GetPostByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        utils.SendError(ctx, http.StatusBadRequest, "Invalid post ID", err)
        return
    }

    post, err := c.postService.GetPostByID(uint(id))
    if err != nil {
        if err.Error() == "record not found" {
            utils.SendError(ctx, http.StatusNotFound, "Post not found", err)
        } else {
            utils.SendError(ctx, http.StatusInternalServerError, "Failed to get post", err)
        }
        return
    }

    utils.SendSuccess(ctx, "Post retrieved successfully", post, nil)
}



func (c *PostController) UploadImage(ctx *gin.Context) {
    file, err := ctx.FormFile("image")
    if err != nil {
        utils.SendError(ctx, http.StatusBadRequest, "Image is required", err)
        return
    }

    allowedMIMETypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
    }

    fileHeader, err := file.Open()
    if err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to open file", err)
        return
    }
    defer fileHeader.Close()

    buffer := make([]byte, 512)
    _, err = fileHeader.Read(buffer)
    if err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to read file header", err)
        return
    }

    contentType := http.DetectContentType(buffer)
    if !allowedMIMETypes[contentType] {
        utils.SendError(ctx, http.StatusBadRequest, "Invalid file content. File is not a valid image format", nil)
        return
    }

    uploadDir := "./uploads"
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, "Failed to create upload directory", err)
        return
    }

    filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
    filepath := filepath.Join(uploadDir, filename)
    
    if err := ctx.SaveUploadedFile(file, filepath); err != nil {
        utils.SendError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to save image to %s", filepath), err)
        return
    }

    utils.SendSuccess(ctx, "Image uploaded successfully", gin.H{"url": "/uploads/" + filename}, nil)
}