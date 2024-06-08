package def

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bytedance/sonic"
)

var (
	day     = 24 * time.Hour
	week    = 7 * day
	quarter = 13 * week
)

// DefaultRecallIntervals 艾宾浩斯记忆曲线的默认复习间隔策略
var DefaultRecallIntervals = RecallIntervalLevel{
	5 * time.Minute,
	30 * time.Minute,
	12 * time.Hour,
	1 * day,
	2 * day,
	4 * day,
	week,
	2 * week,
	4 * week,
	8 * week,
	quarter,
}

// RecallIntervalLevel 定义了复习间隔的级别
type RecallIntervalLevel []time.Duration

// UnmarshalJSON 实现 JSON 反序列化
func (r *RecallIntervalLevel) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*r = DefaultRecallIntervals
		return nil
	}
	var rawIntervals []string
	if err := sonic.Unmarshal(data, &rawIntervals); err != nil {
		return err
	}

	*r = make([]time.Duration, len(rawIntervals))
	for i, rawInterval := range rawIntervals {
		interval, err := time.ParseDuration(rawInterval)
		if err != nil {
			return err
		}
		(*r)[i] = interval
	}

	return nil
}

// MarshalJSON 实现 JSON 序列化
func (r RecallIntervalLevel) MarshalJSON() ([]byte, error) {
	rawIntervals := make([]string, len(r))
	for i, interval := range r {
		rawIntervals[i] = interval.String()
	}

	return sonic.Marshal(rawIntervals)
}

// Scan 实现 sql.Scanner 接口
func (r *RecallIntervalLevel) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return sonic.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口
func (r RecallIntervalLevel) Value() (driver.Value, error) {
	return sonic.Marshal(r)
}

// GetInterval 根据熟练度选择命中哪一个 interval 配置
func (r RecallIntervalLevel) GetInterval(familiarity utils.Percentage) time.Duration {
	ri := DefaultRecallIntervals
	if len(r) != 0 {
		ri = r
	}

	// 动态计算应该使用哪个间隔配置
	index := int(familiarity.Clamp0100()) * len(ri) / 100
	if index >= len(ri) {
		index = len(ri) - 1
	}

	return ri[index]
}
