package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/repositories"
	"github.com/redis/go-redis/v9"
	"github.com/microcosm-cc/bluemonday"
)

type PostService struct {
	repo        *repositories.PostRepository
	redisClient *redis.Client
	likeService *LikeService
}

func NewPostService(repo *repositories.PostRepository, redisClient *redis.Client, likeService *LikeService) *PostService {
    return &PostService{repo: repo, redisClient: redisClient, likeService: likeService}
}

func (s *PostService) CreatePost(post *models.Post) error {

	p := bluemonday.UGCPolicy()  // سیاست امن برای UGC
    post.Content = p.Sanitize(post.Content)  // sanitize قبل از ذخیره
	
	// شروع یک ترنزاکشن جدید
	tx := s.repo.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ذخیره پست در دیتابیس در داخل ترنزاکشن
	if err := s.repo.CreateWithTx(tx, post); err != nil {
		tx.Rollback()
		return err
	}

	// ایجاد الگوی کلید کش برای حذف تمام صفحات مرتبط
	cacheKeyPattern := fmt.Sprintf("posts:lang:%s:type:%s:*", post.Lang, post.Type)
	keys, err := s.redisClient.Keys(context.Background(), cacheKeyPattern).Result()
	if err != nil {
		// اگر خطا در دریافت کلیدها رخ داد، ترنزاکشن را Rollback کن
		tx.Rollback()
		return err
	}
	if len(keys) > 0 {
		if err := s.redisClient.Del(context.Background(), keys...).Err(); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit ترنزاکشن در صورت موفقیت تمام عملیات‌ها
	return tx.Commit().Error
	
}

func (s *PostService) GetPosts(lang string, postType string, page, limit int) ([]models.Post, error) {
	// Generate a unique cache key for the specific query
	cacheKey := fmt.Sprintf("posts:lang:%s:type:%s:page:%d:limit:%d", lang, type, page, limit)
	
	// Try to get posts from cache
	cachedPosts, err := s.redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil && cachedPosts != "" {
		var posts []models.Post
		err = json.Unmarshal([]byte(cachedPosts), &posts)
		if err == nil {
			return posts, nil
		}
	}

	// If cache miss, fetch from database
	posts, err := s.repo.GetByLang(context.Background(), lang, postType, page, limit)  // تغییر به context.Background() اگر ctx نداره
    if err != nil {
        return nil, err
    }

	// اضافه کردن لایک‌ها به هر پست (فرض بر این که userID از context یا request می‌آد، اینجا 0 فرض کردیم برای مثال – در handler واقعی جایگزین کن)
    userID := uint(0)  // جایگزین با userID واقعی از JWT
    for i := range posts {
        likesCount, _ := s.likeService.GetPostLikes(posts[i].ID)
        isLiked, _ := s.likeService.likeRepo.IsLiked(userID, posts[i].ID)
        posts[i].LikesCount = likesCount  // فرض بر این که مدل Post فیلد LikesCount و IsLikedByUser داره (اگر نه، یه struct response جدید بساز)
        posts[i].IsLikedByUser = isLiked
    }

	
	// Serialize posts to JSON and cache them
	serializedPosts, err := json.Marshal(posts)
	if err == nil {
		s.redisClient.Set(context.Background(), cacheKey, serializedPosts, 5*time.Minute) // Cache for 5 minutes
	}

	return posts, nil
}


// تابع جدید برای دریافت یک پست با شناسه
func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
    // ابتدا از کش بررسی می‌کنیم
    cacheKey := fmt.Sprintf("post:%d", id)
    cachedPost, err := s.redisClient.Get(context.Background(), cacheKey).Result()
    if err == nil && cachedPost != "" {
        var post models.Post
        err = json.Unmarshal([]byte(cachedPost), &post)
        if err == nil {
            return &post, nil
        }
    }

    // اگر در کش نبود، از دیتابیس می‌خوانیم
	post, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
	// اضافه کردن لایک‌ها
    userID := uint(0)  // جایگزین با userID واقعی
    likesCount, _ := s.likeService.GetPostLikes(post.ID)
    isLiked, _ := s.likeService.likeRepo.IsLiked(userID, post.ID)
    post.LikesCount = likesCount
    post.IsLikedByUser = isLiked	
	
    // نتیجه را در کش ذخیره می‌کنیم
    postJSON, err := json.Marshal(post)
    if err == nil {
        s.redisClient.Set(context.Background(), cacheKey, postJSON, 1*time.Hour) // کش برای ۱ ساعت
    }

    return post, nil
}