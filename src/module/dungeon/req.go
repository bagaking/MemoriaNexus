package dungeon

import "github.com/bagaking/memorianexus/src/def"

type (
	ReqGetDungeon struct {
		Type def.DungeonType `json:"type"`
	}
)
