package def

import (
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
)

type DungeonType uint8 // 0 ~ 255

const (
	DungeonTypeCampaign DungeonType = 0x1  // 战役地牢
	DungeonTypeEndless  DungeonType = 0x2  // 无尽地牢
	DungeonTypeInstance DungeonType = 0x21 // 即时副本 (随机地牢)
)

func (dt *DungeonType) String() string {
	switch {
	case *dt == DungeonTypeCampaign:
		return "campaign"
	case *dt == DungeonTypeEndless:
		return "endless"
	case *dt >= DungeonTypeInstance:
		return "instance"
	default:
		return "unknown"
	}
}

func (dt *DungeonType) Valid() bool {
	switch *dt {
	case DungeonTypeCampaign:
	case DungeonTypeEndless:
	case DungeonTypeInstance:
	default:
		return false
	}
	return true
}

// UnmarshalJSON custom unmarshaller to handle both strings and numbers
func (dt *DungeonType) UnmarshalJSON(data []byte) error {
	var value any
	if err := sonic.Unmarshal(data, &value); err != nil {
		return err
	}

	switch v := value.(type) {
	case float64:
		*dt = DungeonType(v)
	case string:
		switch strings.TrimSpace(strings.ToLower(v)) {
		case "campaign", "1":
			*dt = DungeonTypeCampaign
		case "endless", "2":
			*dt = DungeonTypeEndless
		case "instance", "3":
			*dt = DungeonTypeInstance
		default:
			return fmt.Errorf("invalid dungeon type: %s", v)
		}
	default:
		return fmt.Errorf("invalid dungeon type: %v", v)
	}

	return nil
}
