// Package handlers provides API request handlers.
package core

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/pkg/auth"
	"github.com/bagaking/memorianexus/src/core/review"
)

// HandleReviewSession 处理复习会话请求
func (svr *Service) HandleReviewSession(c *gin.Context) {
	// 获取认证后的用户ID
	userID, exists := c.Get(auth.UserCtxKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// todo: read newConfidence and reviewTime in request

	// todo: read userFactors and reviewData from DB

	scheduler := review.NewReviewScheduler(userFactors)
	nextReviewTime, newReviewLevel := scheduler.ScheduleNextReview(reviewData, newConfidence)

	// 保存新的复习级别和下一次复习时间到数据库中（伪代码）
	repo.SaveReviewSchedule(c, userID, nextReviewTime, newReviewLevel)

	c.JSON(http.StatusOK, gin.H{
		"next_review_time": nextReviewTime,
		"review_level":     newReviewLevel,
	})
}
