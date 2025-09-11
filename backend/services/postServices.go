package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/repositories"
	"github.com/microcosm-cc/bluemonday"
	"github.com/redis/go-redis/v9"
	"github.com/alimosavifard/zyros-backend/utils"
	"time"
)

// PostResponse: struct واسط برای پاسخ، بدون تغییر مدل Post
type PostResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Type          string `json:"type"`
	Lang          string `json:"lang"`
	ImageUrl      string `json:"imageUrl,omitempty"`
	UserID        uint   `json:"user_id"`
	User          models.User `json:"user,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"` // اگر نیاز باشد
	LikesCount    int64  `json:"likesCount"`
	IsLikedByUser bool   `json:"isLikedByUser"`
}

type PostService struct {
	repo        *repositories.PostRepository
	redisClient *redis.Client
	likeService *LikeService
}

func NewPostService(repo *repositories.PostRepository, redisClient *redis.Client, likeService *LikeService) *PostService {
	return &PostService{repo: repo, redisClient: redisClient, likeService: likeService}
}

func (s *PostService) CreatePost(ctx context.Context, post *models.Post) error {
	p := bluemonday.UGCPolicy()
	post.Content = p.Sanitize(post.Content)

	tx := s.repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := s.repo.CreateWithTx(tx, post); err != nil {
		tx.Rollback()
		return err
	}

	cacheKeyPattern := fmt.Sprintf("posts:lang:%s:type:%s:*", post.Lang, post.Type)
	keys, err := s.redisClient.Keys(ctx, cacheKeyPattern).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(keys) > 0 {
		if err := s.redisClient.Del(ctx, keys...).Err(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *PostService) GetPosts(ctx context.Context, lang string, postType string, page, limit int, userID uint) ([]PostResponse, error) {
	cacheKey := fmt.Sprintf("posts:lang:%s:type:%s:page:%d:limit:%d:user:%d", lang, postType, page, limit, userID)

	cachedPosts, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedPosts != "" {
		var postResponses []PostResponse
		if json.Unmarshal([]byte(cachedPosts), &postResponses) == nil {
			return postResponses, nil
		}
	}

    posts, err := s.repo.GetByLang(ctx, lang, postType, page, limit)
    if err != nil {
        utils.InitLogger().Error().Err(err).Msg("Failed to fetch posts from DB")
        return nil, fmt.Errorf("failed to fetch posts from DB: %w", err)
    }
	

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		likesCount, err := s.likeService.GetPostLikes(ctx, post.ID)
		if err != nil {
			likesCount = 0 // fallback
		}
		isLiked, err := s.likeService.likeRepo.IsLiked(ctx, userID, post.ID)
		if err != nil {
			isLiked = false
		}

		postResponses[i] = PostResponse{
			ID:            post.ID,
			Title:         post.Title,
			Content:       post.Content,
			Type:          post.Type,
			Lang:          post.Lang,
			ImageUrl:      post.ImageUrl,
			UserID:        post.UserID,
			User:          post.User,
			LikesCount:    likesCount,
			IsLikedByUser: isLiked,
		}
	}

	serialized, err := json.Marshal(postResponses)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, serialized, 5*time.Minute)
	}

	return postResponses, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id uint, userID uint) (*PostResponse, error) {
	cacheKey := fmt.Sprintf("post:%d:user:%d", id, userID)

	cachedPost, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedPost != "" {
		var postResp PostResponse
		if json.Unmarshal([]byte(cachedPost), &postResp) == nil {
			return &postResp, nil
		}
	}

	post, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	likesCount, _ := s.likeService.GetPostLikes(ctx, post.ID)
	isLiked, _ := s.likeService.likeRepo.IsLiked(ctx, userID, post.ID)

	postResp := &PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Content:       post.Content,
		Type:          post.Type,
		Lang:          post.Lang,
		ImageUrl:      post.ImageUrl,
		UserID:        post.UserID,
		User:          post.User,
		LikesCount:    likesCount,
		IsLikedByUser: isLiked,
	}

	postJSON, err := json.Marshal(postResp)
	if err == nil {
		s.redisClient.Set(ctx, cacheKey, postJSON, 1*time.Hour)
	}

	return postResp, nil
}