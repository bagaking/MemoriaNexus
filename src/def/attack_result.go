package def

import "github.com/bagaking/memorianexus/internal/utils"

// AttackResult 表示攻击结果的类型
type AttackResult string

const (
	// AttackDefeat 表示攻击失败
	AttackDefeat AttackResult = "defeat"
	// AttackMiss 表示攻击未命中
	AttackMiss AttackResult = "miss"
	// AttackHit 表示攻击命中
	AttackHit AttackResult = "hit"
	// AttackKill 表示攻击致命
	AttackKill AttackResult = "kill"
	// AttackComplete 表示攻击完成
	AttackComplete AttackResult = "complete"
)

// damageRates 定义了每种攻击结果对应的伤害百分比
var damageRates = map[AttackResult]utils.Percentage{
	AttackDefeat:   20,
	AttackMiss:     40,
	AttackHit:      60,
	AttackKill:     80,
	AttackComplete: 100,
}

// DamageRate 返回攻击结果对应的伤害百分比
func (ar AttackResult) DamageRate() utils.Percentage {
	if rate, exists := damageRates[ar]; exists {
		return rate
	}
	return utils.Percentage(0)
}
