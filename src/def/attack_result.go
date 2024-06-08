package def

import "github.com/bagaking/memorianexus/internal/utils"

type (
	AttackResult string
)

const (
	AttackDefeat   AttackResult = "defeat"
	AttackMiss     AttackResult = "miss"
	AttackHit      AttackResult = "hit"
	AttackKill     AttackResult = "kill"
	AttackComplete AttackResult = "complete"
)

func (ar AttackResult) DamageRate() utils.Percentage {
	switch ar {
	case AttackDefeat:
		return utils.Percentage(20)
	case AttackMiss:
		return utils.Percentage(40)
	case AttackHit:
		return utils.Percentage(60)
	case AttackKill:
		return utils.Percentage(80)
	case AttackComplete:
		return utils.Percentage(100)
	}
	return utils.Percentage(0)
}
