package profile

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
)

// ReqUpdateUserSettingsMemorization defines the request format for updating user settings.
type ReqUpdateUserSettingsMemorization struct {
	ReviewIntervalSetting *def.RecallIntervalLevel `json:"review_interval"`
	DifficultyPreference  *utils.Percentage        `json:"difficulty_preference"`
	QuizMode              *string                  `json:"quiz_mode"`
}

// ReqUpdateUserSettingsAdvance defines the request to update advanced settings.
type ReqUpdateUserSettingsAdvance struct {
	Theme              *string `json:"theme"`
	Language           *string `json:"language"`
	EmailNotifications *bool   `json:"email_notifications"`
	PushNotifications  *bool   `json:"push_notifications"`
}

// ReqUpdateProfile defines the request format for the UpdateUserProfile endpoint.
type ReqUpdateProfile struct {
	Nickname  string `json:"nickname,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
}
