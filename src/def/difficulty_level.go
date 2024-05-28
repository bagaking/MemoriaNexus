package def

import (
	"encoding/json"
	"errors"
)

type (
	DifficultyLevel uint8
)

const (
	NoviceNormal    DifficultyLevel = 0x01
	NoviceAdvanced  DifficultyLevel = 0x02
	NoviceChallenge DifficultyLevel = 0x03

	AmateurNormal    DifficultyLevel = 0x11
	AmateurAdvanced  DifficultyLevel = 0x12
	AmateurChallenge DifficultyLevel = 0x13

	ProfessionalNormal    DifficultyLevel = 0x21
	ProfessionalAdvanced  DifficultyLevel = 0x22
	ProfessionalChallenge DifficultyLevel = 0x23

	ExpertNormal    DifficultyLevel = 0x31
	ExpertAdvanced  DifficultyLevel = 0x32
	ExpertChallenge DifficultyLevel = 0x33

	MasterNormal    DifficultyLevel = 0x41
	MasterAdvanced  DifficultyLevel = 0x42
	MasterChallenge DifficultyLevel = 0x43
	MasterExtreme   DifficultyLevel = 0x44
)

var difficultyLevelNames = map[DifficultyLevel]string{
	NoviceNormal:          "novice_normal",
	NoviceAdvanced:        "novice_advanced",
	NoviceChallenge:       "novice_challenge",
	AmateurNormal:         "amateur_normal",
	AmateurAdvanced:       "amateur_advanced",
	AmateurChallenge:      "amateur_challenge",
	ProfessionalNormal:    "professional_normal",
	ProfessionalAdvanced:  "professional_advanced",
	ProfessionalChallenge: "professional_challenge",
	ExpertNormal:          "expert_normal",
	ExpertAdvanced:        "expert_advanced",
	ExpertChallenge:       "expert_challenge",
	MasterNormal:          "master_normal",
	MasterAdvanced:        "master_advanced",
	MasterChallenge:       "master_challenge",
	MasterExtreme:         "master_extreme",
}

func (d *DifficultyLevel) String() string {
	return difficultyLevelNames[*d]
}

// UnmarshalJSON unmarshal the enum from a json string or number
func (d *DifficultyLevel) UnmarshalJSON(data []byte) error {
	var value interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch v := value.(type) {
	case float64:
		*d = DifficultyLevel(uint8(v))
	case string:
		for key, name := range difficultyLevelNames {
			if name == v {
				*d = key
				return nil
			}
		}
		return errors.New("invalid DifficultyLevel name")
	default:
		return errors.New("invalid type for DifficultyLevel")
	}
	return nil
}
