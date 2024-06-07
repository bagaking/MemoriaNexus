package dto

import (
	"time"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/model"
)

// RespGetProfile defines the structure for the user profile API response.
type (
	Profile struct {
		ID        utils.UInt64 `json:"id"`
		Nickname  string       `json:"nickname"`
		Email     string       `json:"email"`
		AvatarURL string       `json:"avatar_url"`
		Bio       string       `json:"bio"`
		CreatedAt time.Time    `json:"created_at"`
		// Include other fields as appropriate.
	}

	SettingsMemorization struct {
		// Definitions should match with ProfileMemorizationSetting
		ReviewIntervalSetting def.RecallIntervalLevel `json:"review_interval_setting"`
		DifficultyPreference  uint8                   `json:"difficulty_preference"`
		QuizMode              string                  `json:"quiz_mode"`
	}

	SettingsAdvance struct {
		Theme              string `json:"theme"`
		Language           string `json:"language"`
		EmailNotifications bool   `json:"email_notifications"`
		PushNotifications  bool   `json:"push_notifications"`
	}

	Points struct {
		Cash     utils.UInt64 `json:"cash"`
		Gem      utils.UInt64 `json:"gem"`
		VIPScore utils.UInt64 `json:"vip_score"`
	}

	RespProfile              = RespSuccess[*Profile]
	RespSettingsMemorization = RespSuccess[*SettingsMemorization]
	RespSettingsAdvance      = RespSuccess[*SettingsAdvance]
	RespPoints               = RespSuccess[*Points]
)

func (p *Profile) FromModel(model *model.Profile) *Profile {
	p.ID = model.ID
	p.Nickname = model.Nickname
	p.Email = model.Email
	p.AvatarURL = model.AvatarURL
	p.Bio = model.Bio
	p.CreatedAt = model.CreatedAt
	return p
}

func (s *SettingsMemorization) FromModel(model *model.ProfileMemorizationSetting) *SettingsMemorization {
	s.ReviewIntervalSetting = model.ReviewIntervalSetting
	s.DifficultyPreference = model.DifficultyPreference
	s.QuizMode = model.QuizMode
	return s
}

func (s *SettingsAdvance) FromModel(model *model.ProfileAdvanceSetting) *SettingsAdvance {
	s.Theme = model.Theme
	s.Language = model.Language
	s.EmailNotifications = model.EmailNotifications
	s.PushNotifications = model.PushNotifications
	return s
}

func (p *Points) FromModel(model *model.ProfilePoints) *Points {
	p.Cash = model.Cash
	p.Gem = model.Gem
	p.VIPScore = model.VipScore
	return p
}
