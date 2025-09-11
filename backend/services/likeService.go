// services/likeService.go

package services

import (
    "github.com/alimosavifard/zyros-backend/models"
    "github.com/alimosavifard/zyros-backend/repositories"
)

type LikeService struct {
    likeRepo *repositories.LikeRepository
}

func NewLikeService(likeRepo *repositories.LikeRepository) *LikeService {
    return &LikeService{likeRepo: likeRepo}
}

func (s *LikeService) LikePost(userID, postID uint) error {
    // Check if the post is already liked
    isLiked, err := s.likeRepo.IsLiked(userID, postID)
    if err != nil {
        return err
    }
    if isLiked {
        return errors.New("post already liked")
    }

    like := &models.PostLike{UserID: userID, PostID: postID}
    return s.likeRepo.Create(like)
}

func (s *LikeService) UnlikePost(userID, postID uint) error {
    return s.likeRepo.Delete(userID, postID)
}

func (s *LikeService) GetPostLikes(postID uint) (int64, error) {
    return s.likeRepo.CountLikes(postID)
}

func (s *LikeService) GetUserLikedPosts(userID uint) ([]models.Post, error) {
    // Implement this function in the repository
    // It would join Post and PostLike tables
    return s.likeRepo.FindPostsByUserID(userID)
}