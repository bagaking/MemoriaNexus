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
		ReviewInterval *def.RecallIntervalLevel `json:"review_interval,omitempty"`
		// 用户挑战偏好
		DifficultyPreference *utils.Percentage `json:"difficulty_preference,omitempty"`
		// 倾向的战斗模式，决定了已经在时间内 monster 出场时，新增和复习的出现策略
		QuizMode *def.QuizMode `json:"quiz_mode,omitempty"`
		// 倾向的战斗模式，决定了已经在时间内 monster 出场时，进行选择的优先级顺序
		PriorityMode *def.PriorityMode `json:"priority_mode,omitempty"`
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

func (s *SettingsMemorization) FromModel(model *model.MemorizationSetting) *SettingsMemorization {
	s.ReviewInterval = &model.ReviewInterval
	s.DifficultyPreference = &model.DifficultyPreference
	s.QuizMode = &model.QuizMode
	s.PriorityMode = &model.PriorityMode
	return s
}

func (s *SettingsMemorization) ToModel(model *model.MemorizationSetting) *model.MemorizationSetting {
	if s == nil {
		return model
	}
	if s.ReviewInterval != nil {
		model.ReviewInterval = *s.ReviewInterval
	}
	if s.DifficultyPreference != nil {
		model.DifficultyPreference = *s.DifficultyPreference
	}
	if s.QuizMode != nil {
		model.QuizMode = *s.QuizMode
	}
	if s.PriorityMode != nil {
		model.PriorityMode = *s.PriorityMode
	}
	return model
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
