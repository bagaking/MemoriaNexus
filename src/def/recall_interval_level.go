package def

import (
	"database/sql/driver"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/bagaking/memorianexus/internal/utils"
)

var (
	minute  = time.Minute
	hour    = time.Hour
	day     = 24 * hour
	week    = 7 * day
	quarter = 13 * week
)

// DefaultRecallIntervals 艾宾浩斯记忆曲线的默认复习间隔策略
var DefaultRecallIntervals = RecallIntervalLevel{
	5 * minute,
	30 * minute,
	12 * hour,
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
// @Description RecallIntervalLevel is a slice of durations in Go duration format (e.g., '1h30m')
type RecallIntervalLevel []time.Duration

// UnmarshalJSON 实现 JSON 反序列化
func (r *RecallIntervalLevel) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*r = DefaultRecallIntervals
		return nil
	}
	var rawIntervals []string
	if err := jsoniter.Unmarshal(data, &rawIntervals); err != nil {
		return err
	}

	*r = make([]time.Duration, len(rawIntervals))
	for i, rawInterval := range rawIntervals {
		interval, err := time.ParseDuration(rawInterval)
		if err != nil {
			return err
		}
		(*r)[i] = time.Duration(interval)
	}

	return nil
}

// MarshalJSON 实现 JSON 序列化
func (r RecallIntervalLevel) MarshalJSON() ([]byte, error) {
	rawIntervals := make([]string, len(r))
	for i, interval := range r {
		rawIntervals[i] = interval.String()
	}

	return jsoniter.Marshal(rawIntervals)
}

// Scan 实现 sql.Scanner 接口
func (r *RecallIntervalLevel) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return jsoniter.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口
func (r RecallIntervalLevel) Value() (driver.Value, error) {
	return jsoniter.Marshal(r)
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
