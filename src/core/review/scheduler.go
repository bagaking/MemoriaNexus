// Package review provides scheduling logic for reviews.

package review

import (
	"time"

	"github.com/bagaking/memorianexus/pkg/memcurve"
)

// ReviewScheduler 提供复习计划调度的功能
type ReviewScheduler struct {
	calculator *memcurve.ReviewCalculator
}

// NewReviewScheduler 创建一个新的复习计划调度器
func NewReviewScheduler(factors memcurve.UserFactors) *ReviewScheduler {
	return &ReviewScheduler{
		// 这里可以根据实际情况初始化适合用户的复习间隔。
		calculator: memcurve.NewReviewCalculator(memcurve.DefaultIntervals, factors),
	}
}

// ScheduleNextReview 使用记忆曲线算法来计算下一次复习的时间
func (s *ReviewScheduler) ScheduleNextReview(data *memcurve.ReviewData, newConfidence float64) (nextReviewTime time.Time, newReviewLevel int) {
	return s.calculator.CalculateNextReview(data, newConfidence)
}
