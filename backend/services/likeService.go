package services

import (
	"context"
	"errors"
	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/repositories"
)

type LikeService struct {
	likeRepo *repositories.LikeRepository
}

func NewLikeService(likeRepo *repositories.LikeRepository) *LikeService {
	return &LikeService{likeRepo: likeRepo}
}

func (s *LikeService) LikePost(ctx context.Context, userID, postID uint) error {
	isLiked, err := s.likeRepo.IsLiked(ctx, userID, postID)
	if err != nil {
		return err
	}
	if isLiked {
		return errors.New("post already liked")
	}

	like := &models.PostLike{UserID: userID, PostID: postID}
	return s.likeRepo.Create(ctx, like)
}

func (s *LikeService) UnlikePost(ctx context.Context, userID, postID uint) error {
	return s.likeRepo.Delete(ctx, userID, postID)
}

func (s *LikeService) GetPostLikes(ctx context.Context, postID uint) (int64, error) {
	return s.likeRepo.CountLikes(ctx, postID)
}

func (s *LikeService) GetUserLikedPosts(ctx context.Context, userID uint) ([]models.Post, error) {
	return s.likeRepo.FindPostsByUserID(ctx, userID)
}